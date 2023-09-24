from datetime import datetime

import pymysql


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
        return result

    def save_detection(self, recording_id, file_name, frame_number, total_frame_number, video_millisecond,
                       total_video_millisecond, confidence):
        with self.connection.cursor() as cursor:
            sql = """
            INSERT INTO detections (recording_id, file_name, frame_number, total_frame_number, video_millisecond, total_video_millisecond, confidence, created_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
            """
            cursor.execute(sql,
                           (recording_id, file_name, frame_number, total_frame_number, video_millisecond,
                            total_video_millisecond, confidence, datetime.utcnow()))
            self.connection.commit()

    def close(self):
        self.connection.close()
