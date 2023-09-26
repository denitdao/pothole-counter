import threading

from flask import Flask, jsonify

from detector import analyze_video

app = Flask(__name__)


@app.route('/')
def status():
    return 'OK'


@app.route('/analyze/<recording_id>', methods=['POST'])
def analyze(recording_id):
    thread = threading.Thread(target=background_analyze, args=(recording_id,))
    thread.start()

    return jsonify({"message": "Analysis started"}), 200


def background_analyze(recording_id):
    ph_number = analyze_video(recording_id)
    print(f"Analysis finished.\nRecording Id: {recording_id}   PH number: {ph_number}")


if __name__ == '__main__':
    app.run()
