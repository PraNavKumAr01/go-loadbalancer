import requests

def send_request():
    """
    Sends a GET request to the load balancer and prints the response.
    """
    url = "http://localhost:8080"
    try:
        response = requests.get(url)
        response.raise_for_status()
        print(response.json())
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    send_request()