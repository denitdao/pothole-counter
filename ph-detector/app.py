import logging
import threading

from flask import Flask, jsonify

from detector import Analyzer

app = Flask(__name__)
logging.basicConfig(level=logging.INFO, format='[%(asctime)s][%(threadName)s] %(levelname)8s: %(message)s')


@app.route('/')
def status():
    return 'OK'


@app.route('/analyze/<recording_id>', methods=['POST'])
def analyze(recording_id):
    thread = threading.Thread(target=background_analyze, args=(recording_id,))
    thread.start()

    return jsonify({"message": "Analysis started"}), 200


def background_analyze(recording_id):
    logging.info(f"Analysis started. \tRecording Id: {recording_id}")
    ph_number = Analyzer().analyze(recording_id)
    logging.info(f"Analysis finished.\tRecording Id: {recording_id}\tPH number: {ph_number}")


if __name__ == '__main__':
    app.run(port=8000)
