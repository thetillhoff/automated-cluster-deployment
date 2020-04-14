FROM debian:stretch

LABEL maintainer="info@enforge.de"


### update os
RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get autoremove -y && \
apt-get clean -y && \
rm -rf /tmp/* /var/tmp/* /var/cache/apt/* /var/cache/distfiles/*


### install required software
## install system basics
RUN apt-get install -y net-tools wget git

## install go to /go/bin/go
RUN wget -q https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz -O /go.tar.gz && \
    tar -zxf /go.tar.gz -C / && \
    rm /go.tar.gz

## go get dependencies
RUN /go/bin/go get gopkg.in/yaml.v3

### startup command
CMD cd /build; /go/bin/go build -o main