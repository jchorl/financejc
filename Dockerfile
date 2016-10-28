FROM golang

RUN mkdir -p /go/src/github.com/jchorl/financejc
WORKDIR /go/src/github.com/jchorl/financejc

RUN openssl genrsa -out server.key 2048 && openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650 -subj '/CN=localhost:443/O=FinanceJC/C=CA'

ADD . /go/src/github.com/jchorl/financejc
RUN go-wrapper download
RUN go-wrapper install

CMD ["./scripts/start.sh"]
