FROM golang

ENV GO111MODULE=on
ENV SRC_DIR=/go/src/github.com/fawind/looper/examples/notes-service/notes-service

WORKDIR /app

ADD . $SRC_DIR
RUN cd $SRC_DIR; go build -o notes-service; cp notes-service /app/

ENTRYPOINT ["./notes-service"]

