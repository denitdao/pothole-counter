from flask import Flask

from detector import analyze_video

app = Flask(__name__)


@app.route('/')
def status():
    return 'OK'


@app.route('/analyze/<video_name>')
def analyze(video_name):
    analyze_video(video_name)
    return 'Analysis started'


if __name__ == '__main__':
    app.run()
