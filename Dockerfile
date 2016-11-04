FROM golang

RUN mkdir -p /go/src/github.com/jchorl/financejc
WORKDIR /go/src/github.com/jchorl/financejc

ADD . /go/src/github.com/jchorl/financejc
RUN go-wrapper download
RUN go-wrapper install
CMD go-wrapper run
