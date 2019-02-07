FROM golang:1.10-alpine AS build

COPY . /go/src/github.com/szabado/zkcli/

RUN apk add curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN cd /go/src/github.com/szabado/zkcli \
  && dep ensure \
  && go install

FROM alpine:3.8

RUN mkdir -p /go/bin
COPY --from=build go/bin/zkcli /go/bin/zkcli
ENV PATH /go/bin:$PATH

ENTRYPOINT ["/go/bin/zkcli"]
