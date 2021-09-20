# syntax=docker/dockerfile:1

FROM golang:1.17.1
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY *.csv ./

RUN go build -o /docker-konomi-project
CMD ["/docker-konomi-project"]