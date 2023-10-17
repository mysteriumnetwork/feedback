FROM golang:1.21-alpine AS builder

# Install packages
RUN apk add --no-cache bash gcc musl-dev linux-headers git curl

# Compile application
WORKDIR /src
ADD . .
RUN go run mage.go -v build

FROM alpine:3.17

# Install application
COPY --from=builder /src/build/feedback /usr/bin/feedback

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD wget localhost:8080/api/v1/ping -q -O - > /dev/null 2>&1

ENTRYPOINT ["/usr/bin/feedback"]
