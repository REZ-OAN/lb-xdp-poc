from flask import Flask
import socket

app = Flask(__name__)

@app.route('/')
def index():
    # Get the server's IP address
    ip_address = socket.gethostbyname(socket.gethostname())
    return f"Served by backend with IP: {ip_address}\n"

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80)
