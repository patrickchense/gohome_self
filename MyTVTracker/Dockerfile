# pulling a lightweight version of golang
FROM golang:1.13-alpine
MAINTAINER Patrick Chen<zhe.chen.sg@gmail.com>
RUN apk --update add --no-cache git

# Copy the local package files to the container's workspace.
ADD . /go/src/MyTVTracker
WORKDIR /go/src/MyTVTracker

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get MyTVTracker

# Run the command by default when the container starts.
ENTRYPOINT ["/go/bin/MyTVTracker"]

# Document that the service listens on port 9000.
EXPOSE 9000