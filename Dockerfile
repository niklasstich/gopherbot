FROM golang:latest AS build
WORKDIR /gopherbot
#prepare depended modules
COPY go.mod .
COPY go.sum .
RUN go mod download

#build app
COPY . .
RUN GOOS=linux go build -a -o gopherbot -v -ldflags "-s -w" .

FROM ubuntu:latest
RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates
COPY --from=build /gopherbot/gopherbot /app/gopherbot
WORKDIR /app
CMD ["/app/gopherbot"]