FROM golang:latest

WORKDIR src/myproject
COPY . .

RUN go get -d -v ./...

RUN pwd && find .
CMD ["go", "run", "main.go"]
