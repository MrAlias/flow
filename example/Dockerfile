FROM golang:1.16-alpine
COPY . /go/src/github.com/MrAlias/flow/
WORKDIR /go/src/github.com/MrAlias/flow/example
RUN go install .
CMD ["/go/bin/example"]
EXPOSE 41820
