FROM golang:1.10.0-alpine

ENV worker_path /go/src/github.com/asoorm/tyk-rmq-middleware/middleware-grpc/
ENV PATH $PATH:$worker_path

WORKDIR $worker_path
COPY . .

RUN go install

ENTRYPOINT middleware-grpc

EXPOSE 9111
