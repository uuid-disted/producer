# Use the Go development container image from devcontainers
FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm

# Set root user to perform privileged actions
USER root

# Copy the entrypoint script and set appropriate permissions
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Set the working directory
WORKDIR /workspaces/manager

# Define the entrypoint script
ENTRYPOINT [ "/usr/local/bin/docker-entrypoint.sh" ]
