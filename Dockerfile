# Dockerfile.production

FROM registry.semaphoreci.com/golang:1.18 as builder

ENV APP_HOME /go/src/heeico

WORKDIR "$APP_HOME"
COPY . $APP_HOME


RUN go mod download
RUN go mod verify
RUN go build -o heeico

FROM registry.semaphoreci.com/golang:1.18

ENV APP_HOME /go/src/heeico
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY client "$APP_HOME"/client
COPY --from=builder "$APP_HOME"/heeico $APP_HOME

EXPOSE 80
CMD ["./heeico"]