# syntax = docker/dockerfile:1-experimental

FROM golang:alpine
RUN mkdir app
ADD . app/
WORKDIR app
COPY go.* .
RUN go mod download
COPY . .
COPY config.yml app/config.yml
RUN go build -o main .
EXPOSE 8080
CMD ["./main", "config.yml"]