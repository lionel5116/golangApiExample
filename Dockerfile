FROM golang:1.17

WORKDIR /go/src/app
COPY . .

RUN mkdir /cmd
RUN go build -o /cmd/app /go/src/app/main.go

ENTRYPOINT ["/cmd/app"]