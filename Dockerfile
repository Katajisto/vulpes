FROM golang:bullseye
COPY . /app
WORKDIR /app

RUN apk update
RUN apk upgrade
RUN GOOS=linux go build -tags prod -o main
WORKDIR /root/
COPY --from=0 /app/main ./
COPY --from=0 /app/views/templates ./views/templates
CMD ["./main"]
EXPOSE 1337
