import logging
import os

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
            # TODO: implement file upload and processing
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
