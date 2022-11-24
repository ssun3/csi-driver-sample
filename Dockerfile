FROM golang:1.15-alpine as dev-env


WORKDIR /app

FROM dev-env as build-env
COPY go.mod /go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /bsos


FROM alpine:3.10 as runtime

COPY --from=build-env /bsos /usr/local/bin/bsos

RUN chmod +x /usr/local/bin/bsos

ENTRYPOINT ["bsos"]

