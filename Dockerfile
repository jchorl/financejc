FROM jchorl/appengine

ADD . /go/src/github.com/jchorl/financejc
WORKDIR /go/src/github.com/jchorl/financejc
RUN go get