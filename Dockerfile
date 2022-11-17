FROM golang:1.18-bullseye AS builder

ENV GOOS=linux 

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN make build

FROM debian:11-slim

WORKDIR /app

RUN apt update -y \
    && apt upgrade -y \
    && apt install -y ca-certificates

COPY --from=builder /usr/src/app/start.sh .
COPY --from=builder /usr/src/app/slashmovie .
COPY --from=chiefy/tired-proxy:1.0.0 /app/tired-proxy .

CMD [ "/app/start.sh" ]