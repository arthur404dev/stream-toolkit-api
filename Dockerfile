FROM golang:latest

WORKDIR /go/app/api

COPY . .

RUN go get ./...
RUN go install ./...

CMD ["api"]