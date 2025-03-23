# Start with an Alpine base image
FROM alpine:3.21

# Create a directory for our application
WORKDIR /app

# Copy the binary, public directory, and .env file
COPY main /app/main

# Ensure the main binary is executable
RUN chmod +x /app/main

# Set the command to run the main binary
CMD ["/app/main"]

