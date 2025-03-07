# Specify the base Docker image. You can read more about
# the available images at https://docs.apify.com/sdk/js/docs/guides/docker-images
# You can also use any other image from Docker Hub.
FROM apify/actor-node:22

# Install system dependencies
RUN apk add --no-cache jq wget bash curl && \
    curl --silent --location \
        https://raw.githubusercontent.com/houseabsolute/ubi/master/bootstrap/bootstrap-ubi.sh | sh

# Install NPM packages, skip optional and development dependencies to
# keep the image small. Avoid logging too much and print the dependency
# tree for debugging
RUN npm --quiet set progress=false \
    && npm install -g apify-cli \
    && echo "Installed NPM packages:" \
    && (npm list --omit=dev --all || true) \
    && echo "Node.js version:" \
    && node --version \
    && echo "NPM version:" \
    && npm --version \
    && rm -r ~/.npm

# Install `freeze`
RUN ubi \
    --verbose \
    --project https://github.com/charmbracelet/freeze \
    --in /usr/local/bin

# Copy run.sh script
COPY ./main.sh ./main.sh

# Run the image.
CMD ["/bin/bash", "main.sh"]
