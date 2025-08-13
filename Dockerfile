FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . /app

RUN go run ./cmd/static-generate/main.go -output ./static

FROM nginx:1.27.5

# copy md and images
COPY --from=builder /app/pages /usr/share/nginx/html

# merge generated html
COPY --from=builder /app/static /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
