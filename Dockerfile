FROM golang

WORKDIR /go/src/randomdogs

ADD . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["randomdogs"]