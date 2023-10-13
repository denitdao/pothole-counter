import pymysql


class Recording:
    def __init__(self, id, file_name, original_file_name, type, status, created_at):
        self.id = id
        self.file_name = file_name
        self.original_file_name = original_file_name
        self.type = type
        self.status = status
        self.created_at = created_at

    @classmethod
    def from_dict(cls, data):
        return cls(data['id'], data['file_name'], data['original_file_name'], data['type'], data['status'],
                   data['created_at'])


class Detection:
    def __init__(self, id, recording_id, file_name, frame_number, total_frame_number,
                 video_millisecond, total_video_millisecond, confidence, created_at):
        self.id = id
        self.recording_id = recording_id
        self.file_name = file_name
        self.frame_number = frame_number
        self.total_frame_number = total_frame_number
        self.video_millisecond = video_millisecond
        self.total_video_millisecond = total_video_millisecond
        self.confidence = confidence
        self.created_at = created_at

    @classmethod
    def from_dict(cls, data):
        return cls(data['id'], data['recording_id'], data['file_name'], data['frame_number'],
                   data['total_frame_number'], data['video_millisecond'], data['total_video_millisecond'],
                   data['confidence'], data['created_at'])


class GPX:
    def __init__(self, id, file_name, status, created_at):
        self.id = id
        self.file_name = file_name
        self.status = status
        self.created_at = created_at

    @classmethod
    def from_dict(cls, data):
        return cls(data['id'], data['file_name'], data['status'], data['created_at'])


class DetectionLocation:
    def __init__(self, id, detection_id, gpx_id, latitude, longitude, created_at):
        self.id = id
        self.detection_id = detection_id
        self.gpx_id = gpx_id
        self.latitude = latitude
        self.longitude = longitude
        self.created_at = created_at

    @classmethod
    def from_dict(cls, data):
        return cls(data['id'], data['detection_id'], data['gpx_id'], data['latitude'], data['longitude'],
                   data['created_at'])


class Database:
    def __init__(self, host, user, password, dbname):
        self.connection = pymysql.connect(
            host=host,
            user=user,
            password=password,
            db=dbname,
            cursorclass=pymysql.cursors.DictCursor
        )

    def get_recording_by_id(self, recording_id):
        with self.connection.cursor() as cursor:
            sql = "SELECT * FROM recordings WHERE id=%s"
            cursor.execute(sql, recording_id)
            result = cursor.fetchone()

        if result:
            return Recording.from_dict(result)
        return None

    def get_gpx_by_id(self, gpx_id):
        with self.connection.cursor() as cursor:
            sql = "SELECT * FROM gpx WHERE id=%s"
            cursor.execute(sql, gpx_id)
            result = cursor.fetchone()

        if result:
            return GPX.from_dict(result)
        return None

    def get_detections_by_recording_id(self, recording_id):
        with self.connection.cursor() as cursor:
            sql = "SELECT * FROM detections WHERE recording_id=%s"
            cursor.execute(sql, recording_id)
            result = cursor.fetchall()

        if result:
            return [Detection.from_dict(item) for item in result]
        return []

    def update_recording_status(self, recording_id, status):
        with self.connection.cursor() as cursor:
            sql = "UPDATE recordings SET status=%s WHERE id=%s"
            cursor.execute(sql, (status, recording_id))
            self.connection.commit()

    def update_gpx_status(self, gpx_id, status):
        with self.connection.cursor() as cursor:
            sql = "UPDATE gpx SET status=%s WHERE id=%s"
            cursor.execute(sql, (status, gpx_id))
            self.connection.commit()

    def save_detection(self, detection: Detection):
        with self.connection.cursor() as cursor:
            sql = """
            INSERT INTO detections (recording_id, file_name, frame_number, total_frame_number, video_millisecond, total_video_millisecond, confidence, created_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
            """
            cursor.execute(sql,
                           (detection.recording_id, detection.file_name, detection.frame_number,
                            detection.total_frame_number, detection.video_millisecond,
                            detection.total_video_millisecond, detection.confidence, detection.created_at))
            self.connection.commit()
            return cursor.lastrowid

    def save_location(self, location: DetectionLocation):
        with self.connection.cursor() as cursor:
            sql = """
            INSERT INTO detection_location (detection_id, gpx_id, latitude, longitude, created_at)
            VALUES (%s, %s, %s, %s, %s)
            """
            cursor.execute(sql, (location.detection_id, location.gpx_id, location.latitude, location.longitude,
                                 location.created_at))
            self.connection.commit()

    def close(self):
        self.connection.close()
