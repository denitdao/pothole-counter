import pymysql


class Recording:
    def __init__(self, id, video_name, original_file_name, status, created_at):
        self.id = id
        self.video_name = video_name
        self.original_file_name = original_file_name
        self.status = status
        self.created_at = created_at

    @classmethod
    def from_dict(cls, data):
        return cls(data['id'], data['video_name'], data['original_file_name'], data['status'], data['created_at'])


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

    def update_recording_status(self, recording_id, status):
        with self.connection.cursor() as cursor:
            sql = "UPDATE recordings SET status=%s WHERE id=%s"
            cursor.execute(sql, (status, recording_id))
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

    def close(self):
        self.connection.close()
