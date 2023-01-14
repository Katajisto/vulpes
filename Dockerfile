FROM golang:bullseye
COPY . /app
WORKDIR /app

RUN apk update
RUN apk upgrade
RUN apk add --update gcc=6.3.0-r4 g++=6.3.0-r4
RUN GOOS=linux go build -tags prod -o main
WORKDIR /root/
COPY --from=0 /app/main ./
COPY --from=0 /app/views/templates ./views/templates
CMD ["./main"]
EXPOSE 1337
