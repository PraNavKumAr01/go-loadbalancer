import requests
import streamlit as st
import subprocess
import os
import random
import time

def send_request(url):
    try:
        response = requests.get(url)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        return {"message": f"Error: {e}"}

def start_server(server_name):
    server_path = os.path.join("servers", f"{server_name}.py")
    subprocess.Popen(["python", server_path])

def start_load_balancer():
    subprocess.Popen(["go", "run", "."])

def main():
    st.title("Load Balancer Demo")

    if 'running_servers' not in st.session_state:
        st.session_state.running_servers = []

    col1, col2, col3 = st.columns(3)
    with col1:
        if st.button("Start Server 1"):
            server_url = "http://localhost:8081"
            if server_url not in st.session_state.running_servers:
                st.session_state.running_servers.append(server_url)
                start_server("server1")
    with col2:
        if st.button("Start Server 2"):
            server_url = "http://localhost:8082"
            if server_url not in st.session_state.running_servers:
                st.session_state.running_servers.append(server_url)
                start_server("server2")
    with col3:
        if st.button("Start Server 3"):
            server_url = "http://localhost:8083"
            if server_url not in st.session_state.running_servers:
                st.session_state.running_servers.append(server_url)
                start_server("server3")

    if st.button("Start Load Balancer", key="start_load_balancer"):
        start_load_balancer()

    num_requests = st.number_input("Enter the number of requests", min_value=1, step=1)

    if st.button("Send Requests"):

        start_time = time.time()

        load_balancer_running = st.session_state.get('load_balancer_running', False)

        if load_balancer_running:
            load_balancer_url = "http://localhost:8080"
            for _ in range(num_requests):
                response = send_request(load_balancer_url)
                st.success(response["message"])
        elif st.session_state.running_servers:
            for _ in range(num_requests):
                server_url = random.choice(st.session_state.running_servers)
                response = send_request(server_url)
                st.success(response["message"])
        else:
            st.warning("No servers running. Start at least one server to send requests.")

        end_time = time.time()

        total_time = end_time - start_time

        with st.sidebar:
            st.info(f"Total time taken : {total_time}")


if __name__ == "__main__":
    main()
