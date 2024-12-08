FROM golang:1.23-bullseye AS builder

ENV GOOS=linux

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN make build

FROM debian:stable-slim

EXPOSE 5000

WORKDIR /app

RUN apt update -y \
    && apt upgrade -y \
    && apt install -y ca-certificates

COPY --from=builder /usr/src/app/slashmovie .

CMD [ "/app/slashmovie" ]