FROM golang:bullseye
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=1 GOOS=linux go build -tags prod -o main

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /app/main ./
COPY --from=0 /app/views/templates ./views/templates
CMD ["./main"]
EXPOSE 1337
