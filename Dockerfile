FROM golang:bullseye
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -tags prod -o main

FROM alpine:latest  
WORKDIR /
COPY --from=0 . .
CMD ["./main"]
EXPOSE 1337
