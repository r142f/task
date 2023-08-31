FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

WORKDIR /backend-trainee-assignment-2023

COPY . .

ENV GOPROXY=direct
RUN go mod download

RUN go build -o /build .

EXPOSE 8080

CMD ["/build"]