FROM golang
WORKDIR /go/src/github.com/jchorl/financejc
ADD . .
RUN go get -v -d
RUN go build -a -o financejc .

FROM golang
COPY --from=0 /go/src/github.com/jchorl/financejc/financejc /go/bin/
CMD financejc
