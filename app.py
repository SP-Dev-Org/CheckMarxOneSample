import os
import sqlite3
from flask import Flask, request, render_template_string
import subprocess
import pickle

app = Flask(__name__)

# SQL Injection vulnerability
@app.route('/user')
def get_user():
    user_id = request.args.get('id')
    conn = sqlite3.connect('database.db')
    cursor = conn.cursor()

    # VULNERABLE: String concatenation in SQL query
    query = "SELECT * FROM users WHERE id = " + user_id
    cursor.execute(query)

    user = cursor.fetchone()
    return f"User: {user}"

# Command Injection vulnerability
@app.route('/execute')
def execute_command():
    filename = request.args.get('file')

    # VULNERABLE: User input in shell command
    os.system(f"cat /var/log/{filename}")

    return "Command executed"

# Path Traversal vulnerability
@app.route('/read')
def read_file():
    filename = request.args.get('filename')

    # VULNERABLE: No path validation
    with open(f"/app/data/{filename}", 'r') as f:
        content = f.read()

    return content

# XSS vulnerability
@app.route('/display')
def display_message():
    message = request.args.get('msg')

    # VULNERABLE: Unescaped user input in HTML
    template = f"<html><body><h1>{message}</h1></body></html>"
    return render_template_string(template)

# Insecure Deserialization
@app.route('/load')
def load_data():
    data = request.args.get('data')

    # VULNERABLE: Deserializing untrusted data
    obj = pickle.loads(data.encode())

    return f"Data loaded: {obj}"

# Hardcoded credentials
def connect_to_api():
    # VULNERABLE: Hardcoded API key
    api_key = "sk_live_1234567890abcdefghijklmnop"
    api_secret = "secret_key_abc123xyz789"

    headers = {
        'Authorization': f'Bearer {api_key}',
        'X-API-Secret': api_secret
    }
    return headers

# SSRF vulnerability
@app.route('/fetch')
def fetch_url():
    url = request.args.get('url')

    # VULNERABLE: Fetching user-provided URL
    import urllib.request
    response = urllib.request.urlopen(url)

    return response.read()

if __name__ == '__main__':
    # VULNERABLE: Debug mode in production
    app.run(debug=True, host='0.0.0.0')
