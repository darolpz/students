FROM golang:alpine

RUN apk update && \
    apk upgrade && \
    apk add git npm

# Set the application directory
WORKDIR $GOPATH/src/github.com/darolpz/students

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

ENV PORT=8080
# Run the executable
CMD ["students"]

