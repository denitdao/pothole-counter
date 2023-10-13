import logging
import os
import xml.etree.ElementTree as ET
from datetime import datetime, timedelta

from GPSPhoto import gpsphoto

from persistence import *
from settings import *


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


class Locator:
    def __init__(self):
        self.db = Database(HOST, USER, PASSWORD, DATABASE)

    def parse_gpx(self, file_path):
        """Parse GPX file and return list of [timestamp, lat, lon]."""
        tree = ET.parse(file_path)
        root = tree.getroot()

        data = []
        for wpt in root.findall('.//wpt'):
            timestamp = wpt.find('./time').text
            lat = wpt.get('lat')
            lon = wpt.get('lon')
            data.append([timestamp, lat, lon])
        return data

    def locate(self, recording_id):
        gpx = self.db.get_gpx_by_id(recording_id)
        if not gpx:
            logging.warning(f"GPX with ID {recording_id} not found!")
            return
        if gpx.status in ["FINISHED", "PROCESSING"]:
            logging.warning(f"GPX with ID {recording_id} already {gpx.status}!")
            return
        self.db.update_gpx_status(gpx.id, "PROCESSING")

        try:
            file_path = os.path.join(STORAGE_PATH, GPX_FOLDER, gpx.file_name)
            location_data = self.parse_gpx(file_path)

            # For demonstration: printing the parsed data
            for item in location_data:
                logging.info(f"Timestamp: {item[0]}, Latitude: {item[1]}, Longitude: {item[2]}")

            detections = self.db.get_detections_by_recording_id(recording_id)
            total_frames = detections[0].total_frame_number  # e.g. 560

            start_time = datetime.fromisoformat(location_data[0][0])
            end_time = datetime.fromisoformat(location_data[-1][0])
            total_duration = (end_time - start_time).total_seconds()  # Total video duration in seconds
            frame_duration = total_duration / total_frames  # Duration of each frame in seconds

            for detection in detections:
                current_frame = detection.frame_number  # e.g. 10

                # Calculate elapsed time for the current frame
                elapsed_time = start_time + timedelta(seconds=frame_duration * current_frame)

                # Find two closest location data points
                prev_point = location_data[0]
                for point in location_data[1:]:
                    if datetime.fromisoformat(point[0]) > elapsed_time:
                        break
                    prev_point = point

                # Interpolate location between prev_point and point based on time
                prev_time = datetime.fromisoformat(prev_point[0])
                next_time = datetime.fromisoformat(point[0])
                fraction = (elapsed_time - prev_time) / (next_time - prev_time)

                lat = float(prev_point[1]) + fraction * (float(point[1]) - float(prev_point[1]))
                lon = float(prev_point[2]) + fraction * (float(point[2]) - float(prev_point[2]))

                # Create and save the DetectionLocation object to the database
                detection_location = DetectionLocation(None, detection.id, gpx.id, lat, lon, datetime.utcnow())
                self.db.save_location(detection_location)

            self.db.update_gpx_status(gpx.id, "FINISHED")
            return True
        except Exception as e:
            logging.error("Error while processing location: " + str(e))
            self.db.update_gpx_status(gpx.id, "FAILED")
            return False
        finally:
            self.db.close()


def get_geolocation(img_path):
    gps_data = gpsphoto.getGPSData(img_path)
    if gps_data is None:
        return None
    logging.info(f"GPS data: {gps_data}")
    return [gps_data['Latitude'], gps_data['Longitude']]
