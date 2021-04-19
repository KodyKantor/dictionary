FROM golang:alpine

WORKDIR /go/src/github.com/kodykantor/dictionary

COPY . .

RUN go install -v ./...

EXPOSE 8080/tcp
CMD ["dictionary", "server"]
