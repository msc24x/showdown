FROM golang:1.20 AS build

ENV APP_DIR /showdown

COPY . $APP_DIR
WORKDIR $APP_DIR

RUN go build cmd/showdown/main.go

FROM msc24x/compilers AS server

RUN adduser showdown \
    && echo "showdown ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

USER showdown

ENV APP_DIR /showdown
WORKDIR $APP_DIR

COPY --from=build $APP_DIR/main server
COPY --from=build $APP_DIR/.env.paths .env.paths

RUN mkdir -p tools/isolate/bin && cp $ISOLATE_PATH tools/isolate/bin/isolate
RUN mkdir tmp

ENV PATH=$PATH$COMPILERS_PATH

CMD [ "sudo", "./server", "-start", "-paths", ".env.paths" ]