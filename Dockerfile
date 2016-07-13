FROM jchorl/appengine

ADD . src/github.com/jchorl/financejc
WORKDIR src/github.com/jchorl/financejc
RUN go get
ENTRYPOINT goapp serve --host=0.0.0.0 appengine