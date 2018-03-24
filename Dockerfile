FROM golang:1.10.0

WORKDIR /app
ADD . /app

RUN go get -v ./...
RUN go build -o bot .

ENV TGBOTTOKEN token
ENV TGPROD 1

CMD ["/app/bot"]