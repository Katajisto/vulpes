FROM golang:alpine
COPY . /app
WORKDIR /app
RUN GOOS=linux go build -tags prod -o main

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /app/main ./
COPY --from=0 /app/views/templates ./views/templates
CMD ["./main"]
EXPOSE 1337
