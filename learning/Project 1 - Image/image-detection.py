import cv2
from ultralytics import YOLO

model = YOLO("../../ph-storage/models/ph_yolov8m.pt")
results = model("11.jpg", device="mps", show=True)

cv2.waitKey(0)
