FROM golang:latest as builder

ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -a -installsuffix cgo -o main .


FROM alpine:latest

RUN apk add --no-cache go

WORKDIR /root/

COPY --from=builder /app/. .

EXPOSE 8080

CMD [ "./main" ]