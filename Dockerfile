FROM golang:1.14 as builder

RUN apt -y update && apt install -y build-essential cmake git && \
  git clone https://github.com/nanomsg/nanomsg.git /tmp/nanomsg && \
  cd /tmp/nanomsg && cmake . && cmake --build . && cmake --build . --target install

ENV APP_USER oraksil
ENV APP_HOME /go/src/app

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

WORKDIR $APP_HOME
USER $APP_USER
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o orakki cmd/app.go


FROM debian:buster

COPY --from=builder /tmp/nanomsg /tmp/nanomsg
RUN apt -y update && apt install -y cmake && \
  cd /tmp/nanomsg && cmake --build . --target install

ENV APP_USER oraksil
ENV APP_HOME /go/src/app
ENV LD_LIBRARY_PATH /usr/local/lib

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY --chown=0:0 --from=builder $APP_HOME/orakki $APP_HOME

USER $APP_USER
CMD ["./orakki"]
