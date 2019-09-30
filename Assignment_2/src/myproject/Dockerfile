FROM golang:latest
#FROM golang:latest AS build


#WORKDIR /Users/nguyen/School/OSWEGO/CSC482/Assignment_2/src/myproject
WORKDIR src/myproject
COPY . .

RUN go get -d -v ./...

#RUN go build -o ./main .

#FROM alpine:latest

#WORKDIR /root/
#COPY --from=build /Users/nguyen/School/OSWEGO/CSC482/Assignment_2/src/myproject ./

RUN pwd && find .
CMD ["go", "run", "main.go"]
