FROM golang:alpine
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -tags prod -o main

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /app/main ./
COPY --from=0 /app/views/templates ./views/templates
COPY --from=0 /app/s3 ./s3
CMD ["./main"]
EXPOSE 1337
