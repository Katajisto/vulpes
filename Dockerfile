FROM golang:bullseye
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=1 GOOS=linux go build -tags prod -o main
CMD ["./main"]
EXPOSE 1337
