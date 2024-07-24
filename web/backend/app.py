from flask import Flask
import os
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)
PORT = os.getenv("PORT")


@app.route("/")
def hello_world():
    return {'status': 'hello world'}




if __name__ == "__main__":
    app.run(host='0.0.0.0', port=PORT, debug=True)
