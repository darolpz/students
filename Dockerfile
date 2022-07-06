FROM golang:alpine

RUN apk update && \
    apk upgrade

# Set the application directory
WORKDIR $GOPATH/src/github.com/darolpz/students

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

#Run swag init to create documents
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.3
RUN swag init

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["students"]

