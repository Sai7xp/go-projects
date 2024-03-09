# Use the official Node.js 18 image as the base for this Docker image
FROM node:18

# Set workdir
WORKDIR /app

# Copy the entrypoint script from the local filesystem into the Docker image
COPY entrypoint.sh /usr/local/bin/

# Make the entrypoint script executable within the Docker image
RUN chmod +x /usr/local/bin/entrypoint.sh

# Specify the entry point for the Docker container to be the entrypoint.sh script
ENTRYPOINT ["entrypoint.sh"]