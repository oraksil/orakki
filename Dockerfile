FROM golang:1.14 as builder

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

ENV APP_USER oraksil
ENV APP_HOME /go/src/app

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY --chown=0:0 --from=builder $APP_HOME/orakki $APP_HOME

USER $APP_USER
CMD ["./orakki"]