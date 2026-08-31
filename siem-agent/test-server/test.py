from flask import Flask, request, jsonify

app = Flask(__name__)
received_logs = []


@app.route("/api/logs", methods=["POST"])
def receive():
    data = request.json
    entries = data.get("entries", [])
    received_logs.append(data)
    print(f"Received {len(entries)} logs from {data.get('source')}")
    return jsonify({"status": "ok"})


@app.route("/api/logs", methods=["GET"])
def get_logs():
    return jsonify({"total": len(received_logs), "recent": received_logs[-10:]})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)