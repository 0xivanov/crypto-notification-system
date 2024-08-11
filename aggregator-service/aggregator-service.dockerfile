FROM golang:1.22-alpine

WORKDIR /app

COPY aggregatorService /app

# Install Delve for debugging
# RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@v1.21.0

# CMD [ "/go/bin/dlv", "--listen=:4000", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/app/aggregatorService"]
CMD ["/app/aggregatorService"]
