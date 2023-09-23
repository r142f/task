FROM golang:latest as builder

ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -a -installsuffix cgo -o main .
RUN go test -c -a -installsuffix cgo -o main_test .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main /app/main_test ./

EXPOSE 8080

CMD [ "./main" ]