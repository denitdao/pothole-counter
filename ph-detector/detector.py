import logging
import math
import os
import time
from datetime import datetime

import cv2
from ultralytics import YOLO

from locator import get_geolocation
from persistence import *
from settings import *
from sort import *

# Constants
MODEL_PATH = os.path.join(STORAGE_PATH, MODEL_FOLDER, MODEL_NAME)
CLASS_NAMES = ["pothole"]
VERTICAL_GAP = 60
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
        self.location = None


class Analyzer:
    def __init__(self):
        self.db = Database(HOST, USER, PASSWORD, DATABASE)
        self.tracker = Sort(max_age=20, min_hits=3, iou_threshold=0.3)
        self.analysis_strategy = {
            "VIDEO": self._analyze_video,
            "IMAGE": self._analyze_image
        }

    def analyze(self, recording_id):
        rec = self.db.get_recording_by_id(recording_id)
        if not rec:
            logging.warning(f"Recording with ID {recording_id} not found!")
            return
        if rec.status in ["FINISHED", "PROCESSING"]:
            logging.warning(f"Recording with ID {recording_id} already {rec.status}!")
            return
        self.db.update_recording_status(rec.id, "PROCESSING")

        try:
            # Prepare state constants
            file_path = os.path.join(STORAGE_PATH, VIDEO_FOLDER, rec.file_name)
            unique_records_folder = rec.file_name.split(".")[0] + "_" + str(int(time.time_ns() / 1_000_000))
            records_path = os.path.join(STORAGE_PATH, RECORD_FOLDER, unique_records_folder)

            unique_results = self.analysis_strategy[rec.type](rec, file_path, unique_records_folder, records_path)

            self.db.update_recording_status(rec.id, "FINISHED")
            return unique_results
        except Exception as e:
            logging.error("Error while processing video: " + str(e))
            self.db.update_recording_status(rec.id, "FAILED")
            return 0
        finally:
            self.db.close()

    def _analyze_video(self, rec, file_path, unique_records_folder, records_path):
        cap = cv2.VideoCapture(file_path)
        video_frames_total = cap.get(cv2.CAP_PROP_FRAME_COUNT)
        video_milliseconds_total = int(cap.get(cv2.CAP_PROP_FRAME_COUNT) * 1_000 / cap.get(cv2.CAP_PROP_FPS))

        state = AnalysisState(rec.id, cap, unique_records_folder, records_path, video_frames_total,
                              video_milliseconds_total)

        while cap.isOpened():
            success, original_image = cap.read()
            if not success:
                break

            # TODO: more clever placement of detection line (use another AI to find horizon on the image?)
            state.detection_line = [0, original_image.shape[0] * 3 // 5, original_image.shape[1],
                                    original_image.shape[0] * 3 // 5]
            state.current_frame = int(cap.get(cv2.CAP_PROP_POS_FRAMES))
            state.current_millisecond = int(cap.get(cv2.CAP_PROP_POS_MSEC))

            # Filling tracker with detected objects at this frame
            detections = collect_detections(original_image)
            tracker_results = self.tracker.update(detections)

            # Storing detected objects
            self._find_records_video(tracker_results, original_image, state)

            print_progress_bar(state.current_frame, state.total_frames, prefix=rec.file_name,
                               suffix="| Detected:" + str(len(state.unique_results)), length=50)

        return len(state.unique_results)

    def _analyze_image(self, rec, file_path, unique_records_folder, records_path):
        original_image = cv2.imread(file_path)

        state = AnalysisState(rec.id, original_image, unique_records_folder, records_path, 0, 0)
        state.location = get_geolocation(file_path)

        # Filling tracker with detected objects
        detections = collect_detections(original_image)
        tracker_results = self.tracker.update(detections)

        # Storing detected objects
        self._find_records_img(tracker_results, original_image, state)

        return len(state.unique_results)

    def _find_records_video(self, records, image, state: AnalysisState):
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
                    self._store_record(record_id, image_copy, confidence, state)

    def _find_records_img(self, records, image, state: AnalysisState):
        for result in records:
            x1, y1, x2, y2 = map(int, result[:4])
            confidence = result[4]
            record_id = int(result[5])

            state.unique_results.append(record_id)
            image_copy = image.copy()
            cv2.rectangle(image_copy, (x1, y1), (x2, y2), (0, 0, 255), 4)  # outline the hole at the frame
            cropped_image = crop_around_detection(image_copy, x1, y1, x2, y2)
            detection_id = self._store_record(record_id, cropped_image, confidence, state)

            if state.location is not None:
                detection_location = DetectionLocation(None, detection_id, None, state.location[0], state.location[1],
                                                       datetime.utcnow())
                self.db.save_location(detection_location)

    def _store_record(self, record_id, image, confidence, state: AnalysisState):
        if not os.path.exists(state.records_path):
            os.makedirs(state.records_path)

        record_file_name = "rec_" + str(record_id) + ".jpg"
        record_file_path = os.path.join(state.records_path, record_file_name)
        cv2.imwrite(record_file_path, image)
        short_path = os.path.join(state.records_folder, record_file_name)

        detection = Detection(None, state.recording_id, short_path, state.current_frame, state.total_frames,
                              state.current_millisecond, state.total_milliseconds, confidence, datetime.utcnow())
        return self.db.save_detection(detection)


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


def crop_around_detection(image, x1, y1, x2, y2):
    h, w, _ = image.shape
    crop_width = w // 2
    crop_height = h // 2

    # Default starting points for cropping
    start_x = 0
    start_y = 0

    # Adjust start_x if detection is on the right half
    if x2 > crop_width:
        start_x = min(w - crop_width, x1)

    # Adjust start_y if detection is on the bottom half
    if y2 > crop_height:
        start_y = min(h - crop_height, y1)

    end_x = start_x + crop_width
    end_y = start_y + crop_height

    cropped_img = image[start_y:end_y, start_x:end_x]
    return cropped_img
