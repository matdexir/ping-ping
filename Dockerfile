FROM golang:1.21-alpine

WORKDIR /app

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /ping-ping

EXPOSE 8080

CMD [ "/ping-ping" ]
