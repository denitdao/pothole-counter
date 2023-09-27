import logging
import math
import os
import time
from datetime import datetime

import cv2
from ultralytics import YOLO

from persistence import *
from settings import *
from sort import *

MODEL_PATH = os.path.join(STORAGE_PATH, MODEL_FOLDER, MODEL_NAME)
CLASS_NAMES = ["pothole"]
VERTICAL_GAP = 25

model = YOLO(MODEL_PATH)


class AnalysisState:
    def __init__(self, recording_id, cap, records_folder, records_path, total_frames, total_milliseconds):
        self.recording_id = recording_id
        self.cap = cap
        self.records_folder = records_folder
        self.records_path = records_path
        self.current_frame = 0
        self.total_frames = total_frames
        self.current_millisecond = 0
        self.total_milliseconds = total_milliseconds
        self.detection_line = []
        self.unique_results = []


# TODO: could refactor to separate processing state and logic
class Analyzer:
    def __init__(self):
        self.db = Database(HOST, USER, PASSWORD, DATABASE)
        self.tracker = Sort(max_age=20, min_hits=3, iou_threshold=0.3)

    def analyze_video(self, recording_id):
        rec = self.db.get_recording_by_id(recording_id)
        if not rec:
            logging.warning(f"Recording with ID {recording_id} not found!")
            return
        if rec.status == "FINISHED":
            logging.warning(f"Recording with ID {recording_id} already analyzed!")
            return
        if rec.status == "PROCESSING":
            logging.warning(f"Recording with ID {recording_id} already processing!")
            return
        self.db.update_recording_status(rec.id, "PROCESSING")

        # Prepare state constants
        video_path = os.path.join(STORAGE_PATH, VIDEO_FOLDER, rec.video_name)
        unique_records_folder = rec.video_name.split(".")[0] + "_" + str(int(time.time_ns() / 1_000_000))
        records_path = os.path.join(STORAGE_PATH, RECORD_FOLDER, unique_records_folder)
        cap = cv2.VideoCapture(video_path)
        video_frames_total = cap.get(cv2.CAP_PROP_FRAME_COUNT)
        video_milliseconds_total = int(cap.get(cv2.CAP_PROP_FRAME_COUNT) * 1_000 / cap.get(cv2.CAP_PROP_FPS))

        state = AnalysisState(rec.id, cap, unique_records_folder, records_path, video_frames_total,
                              video_milliseconds_total)

        while cap.isOpened():
            success, original_image = cap.read()
            if not success:
                break

            # TODO: more clever placement of detection line (use another AI to find horizon on the image?)
            state.detection_line = [0, original_image.shape[0] // 2, original_image.shape[1],
                                    original_image.shape[0] // 2]
            state.current_frame = int(cap.get(cv2.CAP_PROP_POS_FRAMES))
            state.current_millisecond = int(cap.get(cv2.CAP_PROP_POS_MSEC))

            # Filling tracker with detected objects at this frame
            detections = collect_detections(original_image)
            tracker_results = self.tracker.update(detections)

            # Storing detected objects
            self.find_records(tracker_results, original_image, state)

            # Showing progress
            print_progress_bar(state.current_frame, state.total_frames, prefix=rec.video_name,
                               suffix="| Detected:" + str(len(state.unique_results)), length=50)

        self.db.update_recording_status(rec.id, "FINISHED")
        self.db.close()
        return len(state.unique_results)

    def find_records(self, records, image, state: AnalysisState):
        limits = state.detection_line
        for result in records:
            x1, y1, x2, y2 = map(int, result[:4])
            confidence = result[4]
            record_id = int(result[5])
            w, h = x2 - x1, y2 - y1
            center_x, center_y = x1 + w // 2, y1 + h // 2

            if limits[0] < center_x < limits[2] and limits[1] - VERTICAL_GAP < center_y < limits[1] + VERTICAL_GAP:
                if record_id not in state.unique_results:  # Object detected
                    state.unique_results.append(record_id)
                    image_copy = image.copy()
                    cv2.rectangle(image_copy, (x1, y1), (x2, y2), (0, 0, 255), 2)  # outline the hole at the frame
                    self.store_record(record_id, image_copy, confidence, state)

    def store_record(self, record_id, image, confidence, state: AnalysisState):
        if not os.path.exists(state.records_path):
            os.makedirs(state.records_path)

        record_file_name = "rec_" + str(record_id) + ".jpg"
        record_file_path = os.path.join(state.records_path, record_file_name)
        cv2.imwrite(record_file_path, image)
        short_path = os.path.join(state.records_folder, record_file_name)

        detection = Detection(None, state.recording_id, short_path, state.current_frame, state.total_frames,
                              state.current_millisecond, state.total_milliseconds, confidence, datetime.utcnow())
        self.db.save_detection(detection)


def collect_detections(detectable_image):
    results = model(detectable_image, device="mps", stream=True)
    detection_stack = np.empty((0, 5))

    for r in results:
        for box in r.boxes:
            x1, y1, x2, y2 = map(int, box.xyxy[0])
            conf = math.ceil((box.conf[0] * 100)) / 100
            class_name = CLASS_NAMES[int(box.cls[0])]

            # Remembering detected object
            if class_name in ["pothole"] and conf > 0.30:
                detection = np.array([x1, y1, x2, y2, conf])
                detection_stack = np.vstack((detection_stack, detection))

    return detection_stack


def print_progress_bar(iteration, total, prefix='', suffix='', decimals=1, length=100, fill='█'):
    percent = ("{0:." + str(decimals) + "f}").format(100 * (iteration / float(total)))
    filled_length = int(length * iteration // total)
    bar = fill * filled_length + '-' * (length - filled_length)

    logging.info(f'\r{prefix} |{bar}| {percent}% {suffix}')
