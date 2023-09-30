import math

import cv2
import cvzone
from ultralytics import YOLO

from sort import *

CLASS_NAMES = ["pothole"]
DETECTION_LINE = [400, 297, 673, 297]  # TODO: should be dynamic based on the video params
VERTICAL_GAP = 25


def handle_keys():
    key = cv2.waitKey(1) & 0xFF  # milliseconds to wait

    if key == ord('q'):
        cap.release()
        cv2.destroyAllWindows()
    elif key == ord('p'):
        while True:
            key2 = cv2.waitKey(1) & 0xFF
            if key2 == ord('q'):
                cap.release()
                cv2.destroyAllWindows()
                break
            elif key2 == ord('p'):
                break


def process_detections(image, mask):
    # detectable_image = cv2.bitwise_and(image, mask)
    detectable_image = image
    results = model(detectable_image, device="mps", stream=True)
    detection_stack = np.empty((0, 5))

    for r in results:
        for box in r.boxes:
            # Boundaries of the object
            x1, y1, x2, y2 = map(int, box.xyxy[0])
            # Confidence
            conf = math.ceil((box.conf[0] * 100)) / 100
            # Class Name
            class_name = CLASS_NAMES[int(box.cls[0])]

            # Remembering detected object
            if class_name in ["pothole"] and conf > 0.30:
                detection = np.array([x1, y1, x2, y2, conf])
                detection_stack = np.vstack((detection_stack, detection))

    return detection_stack


def draw_objects(image, results, limits, unique_results):
    cv2.line(image, (limits[0], limits[1]), (limits[2], limits[3]), (0, 0, 255), 2)

    for result in results:
        x1, y1, x2, y2, result_id = map(int, result)
        w, h = x2 - x1, y2 - y1
        # highlighting object
        cvzone.cornerRect(image, (x1, y1, w, h), l=5, t=2, rt=2, colorR=(255, 0, 0))
        cvzone.putTextRect(image, f'id: {int(result_id)}', (max(0, x1), max(35, y1)),
                           scale=1, thickness=2, offset=5, colorR=(255, 0, 0))

        center_x, center_y = x1 + w // 2, y1 + h // 2
        cv2.circle(image, (center_x, center_y), 2, (0, 255, 0), cv2.FILLED)

        if limits[0] < center_x < limits[2] and limits[1] - VERTICAL_GAP < center_y < limits[1] + VERTICAL_GAP:
            if unique_results.count(result_id) == 0:
                # Object detected. TODO: image of the detection could be saved with metadata and location.
                unique_results.append(result_id)
                cv2.line(image, (limits[0], limits[1]), (limits[2], limits[3]), (0, 255, 0), 4)

    cv2.putText(image, str(len(unique_results)), (255, 100), cv2.FONT_HERSHEY_PLAIN, 5, (50, 50, 255), 8)

    return image


if __name__ == '__main__':
    cap = cv2.VideoCapture("../../ph-storage/videos/potholes.mp4")
    mask = cv2.imread("mask.png")
    graphics = cv2.imread("graphics.png", cv2.IMREAD_UNCHANGED)

    # Detecting and Tracking parts
    model = YOLO("../../ph-storage/models/ph_yolov8n.pt")
    tracker = Sort(max_age=20, min_hits=3, iou_threshold=0.3)

    total_count = []

    while cap.isOpened():
        success, original_image = cap.read()
        if not success:
            break
        DETECTION_LINE = [0, original_image.shape[0] * 3 // 5, original_image.shape[1], original_image.shape[0] * 3 // 5]

        # Filling tracker with detected objects at this frame
        detections = process_detections(original_image, mask)
        tracker_results = tracker.update(detections)

        # Drawing
        # final_image = cvzone.overlayPNG(original_image, graphics, (0, 0))
        final_image = original_image
        final_image = draw_objects(final_image, tracker_results, DETECTION_LINE, total_count)

        cv2.imshow("Image", final_image)

        handle_keys()

    print("Total Count: ", len(total_count))
    # TODO: save all recorded potholes with metadata to the database.
    #  in the UI allow to remove, look at the image, etc.


# OUTSTANDING QUESTIONS:
# How to use custom model for pothole detection? +
# How to run this script as some service?
# How to save all detected potholes to the database?
# Optimization: how to set correct detection line and tune to avoid potholes that are too close to each other?