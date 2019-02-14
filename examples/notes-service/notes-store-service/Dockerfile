FROM golang

ENV GO111MODULE=on
ENV SRC_DIR=/go/src/github.com/fawind/docker-test/examples/notes-service/notes-store-service

WORKDIR /app

ADD . $SRC_DIR
RUN cd $SRC_DIR; go build -o notes-store-service; cp notes-store-service /app/

ENTRYPOINT ["./notes-store-service"]
