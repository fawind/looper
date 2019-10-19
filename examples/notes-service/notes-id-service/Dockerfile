FROM golang

ENV GO111MODULE=on
ENV SRC_DIR=/go/src/github.com/fawind/looper/examples/notes-service/notes-id-service

WORKDIR /app

ADD . $SRC_DIR
RUN cd $SRC_DIR; go build -o notes-id-service; cp notes-id-service /app/

ENTRYPOINT ["./notes-id-service"]
