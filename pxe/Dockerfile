FROM debian:bookworm

LABEL maintainer="info@enforge.de"


### update os
RUN apt-get update && \
apt-get upgrade -y && \
apt-get autoremove -y && \
apt-get clean -y && \
rm -rf /tmp/* /var/tmp/* /var/cache/apt/* /var/cache/distfiles/*


### install required software
## install system basics
RUN apt-get install -y net-tools
## install dnsmasq
RUN apt-get install -y dnsmasq
## install nginx
RUN apt-get install -y nginx

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

### add custom files
COPY ./start.sh /start.sh
RUN chmod +x /start.sh
RUN mkdir /key /etc/ansible


### startup command
CMD /start.sh && service nginx start && dnsmasq -C /etc/dnsmasq.conf -u root && cd /container/go/; ./main