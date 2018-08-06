FROM golang:1.10.0-alpine

ENV worker_path /go/src/github.com/asoorm/tyk-rmq-middleware/worker/
ENV PATH $PATH:$worker_path

WORKDIR $worker_path
COPY . .

RUN go install

ENTRYPOINT worker
