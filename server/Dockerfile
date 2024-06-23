FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

RUN apt-get update && \
    apt-get install -y iproute2 iputils-ping tcpdump curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
# Install Flask in the container
RUN pip install flask

# Make port 80 available to the world outside this container
EXPOSE 80

# Define environment variable
ENV NAME World

# Run app.py when the container launches
CMD ["python", "app.py"]
