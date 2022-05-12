# Dockerfile.production

FROM registry.semaphoreci.com/golang:1.18 as builder

ENV APP_HOME /go/src/heeico

WORKDIR "$APP_HOME"
COPY src/ .

RUN go mod download
RUN go mod verify
RUN go build -o heeico

FROM registry.semaphoreci.com/golang:1.18

ENV APP_HOME /go/src/heeico
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY src/conf/ conf/
COPY src/views/ views/
COPY --from=builder "$APP_HOME"/heeico $APP_HOME

EXPOSE 8080
CMD ["./heeico"]