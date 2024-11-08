FROM golang:1.22.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -a -o /app/bin/fusion-solar-scrapper . \
    && chmod +x /app/bin/fusion-solar-scrapper \
    && cp /app/bin/fusion-solar-scrapper /usr/local/bin/fusion-solar-scrapper

FROM --platform=linux/amd64 node:lts-bullseye-slim
ARG DEBIAN_FRONTEND=noninteractive

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ADD https://go.dev/dl/go1.22.0.linux-amd64.tar.gz /tmp/go.tar.gz

RUN apt-get update \
    && apt-get install --no-install-recommends -y build-essential ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$PATH

RUN npx -y playwright@1.48.1 install --with-deps
COPY --from=builder /app/bin/fusion-solar-scrapper .

ENTRYPOINT ["/fusion-solar-scrapper"]
