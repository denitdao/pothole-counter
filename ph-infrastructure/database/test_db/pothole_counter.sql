DROP DATABASE IF EXISTS pothole_counter;
CREATE DATABASE IF NOT EXISTS pothole_counter;
USE pothole_counter;

SELECT 'CREATING DATABASE STRUCTURE' as 'INFO';

DROP TABLE IF EXISTS recordings,
    detections,
    gpx,
    detection_location;

/*!50503 set default_storage_engine = InnoDB */;
/*!50503
select CONCAT('storage engine: ', @@default_storage_engine) as INFO */;

CREATE TABLE recordings
(
    id                 INT          NOT NULL AUTO_INCREMENT,

    video_name         VARCHAR(255) NOT NULL,
    original_file_name VARCHAR(255) NOT NULL,
    created_at         TIMESTAMP    NOT NULL,

    UNIQUE (video_name),
    PRIMARY KEY (id)
);

CREATE TABLE detections
(
    id                INT          NOT NULL AUTO_INCREMENT,
    recording_id      INT          NOT NULL,

    file_name         VARCHAR(255) NOT NULL,
    frame_number      INT          NOT NULL,
    video_millisecond INT          NOT NULL,
    confidence        FLOAT        NOT NULL,
    created_at        TIMESTAMP    NOT NULL,

    FOREIGN KEY (recording_id) REFERENCES recordings (id) ON DELETE CASCADE,
    PRIMARY KEY (id, recording_id)
);

CREATE TABLE gpx
(
    id           INT          NOT NULL AUTO_INCREMENT,
    recording_id INT          NOT NULL,

    file_name    VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP    NOT NULL,

    FOREIGN KEY (recording_id) REFERENCES recordings (id) ON DELETE CASCADE,
    PRIMARY KEY (id, recording_id)
);

CREATE TABLE detection_location
(
    id           INT       NOT NULL AUTO_INCREMENT,
    detection_id INT       NOT NULL,
    gpx_id       INT       NOT NULL,

    latitude     FLOAT     NOT NULL,
    longitude    FLOAT     NOT NULL,
    created_at   TIMESTAMP NOT NULL,

    FOREIGN KEY (detection_id) REFERENCES detections (id) ON DELETE CASCADE,
    FOREIGN KEY (gpx_id) REFERENCES gpx (id) ON DELETE CASCADE,
    PRIMARY KEY (id, detection_id, gpx_id)
);

flush /*!50503 binary */ logs;

SELECT 'LOADING recordings' as 'INFO';
source docker-entrypoint-initdb.d/load_recordings.dump ;
SELECT 'LOADING detections' as 'INFO';
source docker-entrypoint-initdb.d/load_detections.dump ;
SELECT 'LOADING gpx' as 'INFO';
source docker-entrypoint-initdb.d/load_gpx.dump ;
SELECT 'LOADING dept_manager' as 'INFO';
source docker-entrypoint-initdb.d/load_detection_location.dump ;
