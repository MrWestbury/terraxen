FROM golang:1.17-alpine AS build
LABEL maintainer="Terraxen Authors <developers@wifu.co.uk>"

WORKDIR /app

COPY cmd/       /app/cmd
COPY backend/   /app/backend
COPY services/  /app/services

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o /goapps/env-man /app/cmd/environment-manager/main.go

FROM golang:1.17-alpine

RUN adduser --disabled-password --gecos '' noroot

COPY --from=build /goapps/* /app

USER noroot:noroot

EXPOSE 8000

ENTRYPOINT ["/app"]

CMD ["env-man"]
