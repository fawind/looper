FROM alpine

RUN apk add --no-cache curl

COPY . /app
WORKDIR /app

ENTRYPOINT ["sh", "run.sh"]

