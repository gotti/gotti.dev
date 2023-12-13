FROM golang:1.20-alpine

WORKDIR /app
COPY . /app

RUN go build -o main ./app/app.go

CMD ["/app/main"]
