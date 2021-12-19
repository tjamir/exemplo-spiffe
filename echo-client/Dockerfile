FROM golang:1.17-alpine as builder
ARG APP_NAME=echo-client
RUN apk add build-base git mercurial
ADD . /opt/$APP_NAME
RUN cd /opt/$APP_NAME && go mod tidy && go mod download
WORKDIR /opt/$APP_NAME
RUN go build -o /opt/$APP_NAME/bin/$APP_NAME .

FROM spire-agent:latest-local as agent

FROM alpine as echo-client
ARG APP_NAME=echo-client
ENV APP_NAME_ENV=$APP_NAME
COPY --from=builder /opt/$APP_NAME/bin/$APP_NAME /opt/$APP_NAME/bin/$APP_NAME
COPY --from=agent /opt/spire /opt/spire
WORKDIR /opt/$APP_NAME
ENTRYPOINT /opt/$APP_NAME_ENV/bin/$APP_NAME_ENV