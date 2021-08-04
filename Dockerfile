FROM golang:latest AS build
WORKDIR /gopherbot
COPY . .
RUN go get -d -v .
RUN GOOS=linux go build -a -o gopherbot -v -ldflags "-s -w" .

FROM ubuntu:latest
RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates
COPY --from=build /gopherbot/gopherbot /app/gopherbot
COPY --from=build /gopherbot/config/config.toml /app/config/config.toml
WORKDIR /app
CMD ["/app/gopherbot"]