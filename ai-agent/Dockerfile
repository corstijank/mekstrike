FROM python:3.11-slim

WORKDIR /app

# Copy requirements and install dependencies
COPY ai-agent/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY ai-agent/ .

# Expose ports
EXPOSE 50051 8000

# Run the application
CMD ["python", "app.py"]