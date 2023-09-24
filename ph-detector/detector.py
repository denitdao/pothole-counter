import math

import cv2
from ultralytics import YOLO

from persistence import *
from settings import *
from sort import *

MODEL_PATH = os.path.join(STORAGE_PATH, MODEL_FOLDER, MODEL_NAME)
CLASS_NAMES = ["pothole"]
VERTICAL_GAP = 25

model = YOLO(MODEL_PATH)
db = Database(HOST, USER, PASSWORD, DATABASE)


def analyze_video(recording_id):
    recording = db.get_recording_by_id(recording_id)
    if not recording:
        print(f"Recording with ID {recording_id} not found!")
        return

    video_name = recording['video_name']
    video_path = os.path.join(STORAGE_PATH, VIDEO_FOLDER, video_name)
    unique_video_folder = video_name.split(".")[0] + "_" + str(int(time.time_ns() / 1_000_000))
    cap = cv2.VideoCapture(video_path)
    video_frames_total = cap.get(cv2.CAP_PROP_FRAME_COUNT)

    # Detecting and Tracking parts
    tracker = Sort(max_age=20, min_hits=3, iou_threshold=0.3)
    total_count = []

    while cap.isOpened():
        success, original_image = cap.read()
        if not success:
            break

        # Filling tracker with detected objects at this frame
        detections = collect_detections(original_image)
        tracker_results = tracker.update(detections)

        # Drawing
        detection_line = [0, original_image.shape[0] // 2, original_image.shape[1], original_image.shape[0] // 2]
        total_count = record_objects(tracker_results, detection_line, total_count, original_image, unique_video_folder,
                                     cap, recording_id)

        video_frame = cap.get(cv2.CAP_PROP_POS_FRAMES)
        print_progress_bar(video_frame, video_frames_total, prefix=video_name,
                           suffix="| Detected:" + str(len(total_count)), length=50)

    return len(total_count)


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


def record_objects(results, limits, unique_results, image, folder_name, cap, recording_id):
    for result in results:
        x1, y1, x2, y2, result_id = map(int, result)
        w, h = x2 - x1, y2 - y1
        center_x, center_y = x1 + w // 2, y1 + h // 2

        if limits[0] < center_x < limits[2] and limits[1] - VERTICAL_GAP < center_y < limits[1] + VERTICAL_GAP:
            if unique_results.count(result_id) == 0:
                # Object detected
                unique_results.append(result_id)
                # outline the hole at the frame and store locally
                cv2.rectangle(image, (x1, y1), (x2, y2), (0, 0, 255), 2)
                store_record(folder_name, result_id, image, cap, recording_id)

    return unique_results


def print_progress_bar(iteration, total, prefix='', suffix='', decimals=1, length=100, fill='█'):
    percent = ("{0:." + str(decimals) + "f}").format(100 * (iteration / float(total)))
    filled_length = int(length * iteration // total)
    bar = fill * filled_length + '-' * (length - filled_length)

    print(f'\r{prefix} |{bar}| {percent}% {suffix}')


def store_record(folder, record_id, image, cap, recording_id):
    record_path = os.path.join(STORAGE_PATH, RECORD_FOLDER, folder)
    if not os.path.exists(record_path):
        os.makedirs(record_path)

    record_file_name = "rec_" + str(record_id) + ".jpg"
    record_file = os.path.join(record_path, record_file_name)
    cv2.imwrite(record_file, image)
    print(record_file)

    frame_number = int(cap.get(cv2.CAP_PROP_POS_FRAMES))
    file_name = os.path.join(folder, record_file_name)
    total_frame_number = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
    video_millisecond = int(cap.get(cv2.CAP_PROP_POS_MSEC))
    total_video_millisecond = int(cap.get(cv2.CAP_PROP_FRAME_COUNT) * 1000 / cap.get(cv2.CAP_PROP_FPS))
    confidence = 0.95
    db.save_detection(recording_id, file_name, frame_number, total_frame_number, video_millisecond,
                      total_video_millisecond, confidence)
