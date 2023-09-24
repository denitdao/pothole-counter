from flask import Flask

from detector import analyze_video

app = Flask(__name__)


@app.route('/')
def status():
    return 'OK'


@app.route('/analyze/<recording_id>')
def analyze(recording_id):
    analyze_video(recording_id)
    return 'Analysis started'


if __name__ == '__main__':
    app.run()
