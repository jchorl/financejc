FROM jchorl/appengine

ADD . /go/src/github.com/jchorl/financejc
WORKDIR /go/src/github.com/jchorl/financejc
RUN go get
ENTRYPOINT /go_appengine/goapp serve --host=0.0.0.0 appengine