import requests
from concurrent.futures import ThreadPoolExecutor, as_completed
import random
import time


# Persistent session
session = requests.Session()

# Define the URL for the service
service_url = "http://server-service.default.svc.cluster.local/api"

# for i in range(10000):
#     response = session.get(service_url)
#     print(response.text)

# Function to send request and return the response
def send_request(i):
    try:
        response = session.get(service_url, timeout=5)  # Adjust timeout as necessary
        return f"Request {i} - {response.text}"
    except requests.exceptions.Timeout:
        return f"Request {i} timed out."
    except requests.exceptions.RequestException as e:
        return f"Request {i} failed: {str(e)}"

# Using ThreadPoolExecutor to send 10000 requests concurrently
with ThreadPoolExecutor(max_workers=10) as executor:
    futures = [executor.submit(send_request, i) for i in range(100000)]

    # Wait for each request to complete and print the result as it comes
    for future in as_completed(futures):
        print(future.result())


# Keep the script alive to prevent the pod from restarting
while True:
    time.sleep(3600)  # Sleep for an hour or more to keep the pod running.
