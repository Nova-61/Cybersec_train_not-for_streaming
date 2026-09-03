FROM python:3.11-alpine
RUN pip install --no-cache-dir flask
WORKDIR /app
COPY app.py .
CMD ["python3", "app.py"]