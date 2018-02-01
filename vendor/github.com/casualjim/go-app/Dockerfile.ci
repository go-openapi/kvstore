FROM golang:1.7

RUN apt-get update -yqq &&\
  apt-get install -yqq haveged rsyslog gnupg2 &&\
  go get -u github.com/axw/gocov/gocov &&\
  go get -u gopkg.in/matm/v1/gocov-html &&\
  go get -u github.com/cee-dub/go-junit-report
