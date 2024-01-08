# Pothole Counter

A web application that performs video and image analysis, detecting potholes with computer vision.

## Web application pages

All recordings view:

<img width="600" alt="image" src="https://github.com/denitdao/pothole-counter/assets/49095078/b93f7317-953a-4e3d-88ba-64ac73b6f75d">

Video recording analysis results:

<img width="600" alt="image" src="https://github.com/denitdao/pothole-counter/assets/49095078/0ce07597-a477-400c-ba8b-c635d2ccd10f">

Image recording analysis results:

<img width="600" alt="image" src="https://github.com/denitdao/pothole-counter/assets/49095078/7ccffe8a-3df9-42b1-9296-3d043ce77300">

Map view:

<img width="600" alt="image" src="https://github.com/denitdao/pothole-counter/assets/49095078/d0f1a276-3cca-41f2-8191-2fdd895ee705">

# Technical details

## Technologies

1. `Go` - web application and processing management
2. `Gin` - go-router, handling url paths and templates
3. `HTMX`, `Tailwind` - interactive webpages
4. `AlpineJS` - Google Maps data management
5. `Python`, `Flask` - web application for video/image analysis
6. `YOLOv8` - video analysis AI model
7. `MySQL` - data storage
8. `Docker` - containerizing DB

## Structure

Application contains 2 main modules.

`ph-manager` - webserver on Golang, hosting web application (HTMX) and managing creation and deletion of the recordings.

`ph-detector` - webserver on Python running a YOLOv8 AI model and processing video/image data to discover and store potholes and find their geolocation using GPX files. 

`ph-storage` - defines the filesystem for this project, storing Videos, Images, Detections, GPX files and Yolo Models.

The main processing flows look like this:

<img width="600" alt="image" src="https://github.com/denitdao/pothole-counter/assets/49095078/c7409bee-1582-4b84-a62a-377f06d88435">

## DB

Application uses MySQL database to store data about recordings (uploaded images, videos) and detection (images of the potholes discovered)

<img width="600" alt="image" src="https://github.com/denitdao/pothole-counter/assets/49095078/8ecefb9d-35c7-4fc9-bd92-e513a72c84f3">


