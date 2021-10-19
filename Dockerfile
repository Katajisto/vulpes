FROM golang:alpine
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -tags prod -o main

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /app/main ./
CMD ["./main"]
