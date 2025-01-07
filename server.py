from flask import Flask, request
import os
import random
import time

app = Flask(__name__)

@app.route('/api', methods=['GET'])
def api():
    delay = random.randint(1, 3)
    pod_id = os.getenv("POD_ID", "unknown")  # Use an environment variable to identify the Pod
    print(f"Pod {pod_id} will delay the response for {delay} seconds.")
    time.sleep(delay)
    return f"Response from Pod {pod_id}, delayed by {delay} seconds", 200

if __name__ == "__main__":
    pod_id = os.getenv("POD_ID", "default")
    app.run(host="0.0.0.0", port=5000)

