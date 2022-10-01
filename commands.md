# Commands to build and run cardWithWords application

All commands should be run from the root directory.

## Build docker container

```shell script
# Build the application in the root directory (with `go.mod` file)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/cardWithWords ./cmd/cardWithWords

# cd to docker dir and create docker container 
docker build -t cardWithWords:latest -f ./Dockerfile .

# for mac m1
GOARCH=amd64 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ./docker/cardWithWords ./cmd/cardWithWords
docker buildx build --platform linux/amd64 -o type=docker -t cardwithwords:latest -f ./docker/Dockerfile ./docker
```

## Run docker container

```shell script
# Run docker container
docker run -d --name cardWithWords --restart always cardwithwords:latest /
cardWithWords --token YOUR_TOKEN --wordsQuantity QUANTITY
```