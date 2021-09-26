FROM golang:1.17.1-buster AS build

COPY ../ /build/

WORKDIR /build

RUN go build -o /build/out/ .

FROM build

COPY /build/out/ /app

ENTRYPOINT ["/app/terrakube-moduleregistry"]

