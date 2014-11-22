FROM golang:1.3.3
MAINTAINER Bayu Aldi Yansyah <bayualdiyansyah@gmail.com>

ENV PORT 8080
EXPOSE 8080

RUN go get github.com/lib/pq

ADD . /go/src/github.com/pyk/automata
RUN go install github.com/pyk/automata