import cv2
import math
from flask import Flask
from ultralytics import YOLO

from sort import *

basepath = os.path.dirname(__file__)
videopath = os.path.join(basepath, "../../ph-storage/videos/potholes.mp4")
modelpath = os.path.join(basepath, "../../ph-storage/models/ph_yolov8n.pt")

CLASS_NAMES = ["pothole"]
VERTICAL_GAP = 25
model = YOLO(modelpath)

app = Flask(__name__)


@app.route("/")
def hello_world():
    cap = cv2.VideoCapture(videopath)
    total_frames = cap.get(cv2.CAP_PROP_FRAME_COUNT)

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
        total_count = record_objects(tracker_results, detection_line, total_count)

        frame_number = cap.get(cv2.CAP_PROP_POS_FRAMES)
        process_percent = math.ceil((frame_number / total_frames) * 100)
        print("Processed: ", process_percent, "%")
        print("Detected:  ", len(total_count))

    return "Total Count: " + str(len(total_count))


def collect_detections(detectable_image, mask=None):
    if mask is not None:
        detectable_image = cv2.bitwise_and(detectable_image, mask)

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


def record_objects(results, limits, unique_results):
    for result in results:
        x1, y1, x2, y2, result_id = map(int, result)
        w, h = x2 - x1, y2 - y1
        center_x, center_y = x1 + w // 2, y1 + h // 2

        if limits[0] < center_x < limits[2] and limits[1] - VERTICAL_GAP < center_y < limits[1] + VERTICAL_GAP:
            if unique_results.count(result_id) == 0:
                # Object detected
                unique_results.append(result_id)

    return unique_results
