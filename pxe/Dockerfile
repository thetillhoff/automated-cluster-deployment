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
RUN apt-get install -y net-tools wget
## install dnsmasq
RUN apt-get install -y dnsmasq
## install nginx
RUN apt-get install -y nginx
## install golang
RUN wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz -O /go.tar.gz && \
  mkdir -p /usr/local/go && tar -zxf /go.tar.gz -C /usr/local/ && rm /go.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

## add nginx config file
COPY ./nginx_site_config /etc/nginx/sites-available/site
## enable nginx site and remove default site
RUN rm /etc/nginx/sites-enabled/default && ln -s /etc/nginx/sites-available/site /etc/nginx/sites-enabled/
EXPOSE 80/tcp

## add dnsmasq config file
COPY ./dnsmasq.conf /etc/dnsmasq.conf
# 53 dns
# 69/udp tftp
EXPOSE 53/udp
EXPOSE 69/udp
EXPOSE 69/tcp

## add go files
COPY golang /go


### add custom files
COPY ./start.sh /start.sh
RUN chmod +x /start.sh
RUN mkdir /key /etc/ansible


### startup command
CMD /start.sh && service nginx start && dnsmasq && go build /go/ -o /go/http && /go/http