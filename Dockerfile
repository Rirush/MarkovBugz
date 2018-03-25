FROM golang:1.10.0

ENV GOBIN /go/bin

RUN mkdir /app
RUN mkdir /go/src/app

ADD . /go/src/app
WORKDIR /go/src/app

RUN go get -v ./...
RUN go build -o /app/bot .

ENV TGBOTTOKEN token
ENV TGPROD 1
CMD ["/app/bot"]