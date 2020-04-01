FROM debian:stretch

LABEL maintainer="info@enforge.de"


### update and install os
RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get autoremove && \
apt-get clean -y && \
rm -rf /tmp/* /var/tmp/* /var/cache/apt/* /var/cache/distfiles/*


### install required software
## install system basics
RUN apt-get install -y net-tools
## install nginx with php
RUN apt-get install -y nginx
RUN apt-get install -y php7.0-fpm
## install ansible
RUN apt-get install -y gnupg2 curl && \
echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" >> /etc/apt/sources.list && \
curl -sL "http://keyserver.ubuntu.com/pks/lookup?op=get&search=0x93C4A3FD7BB9C367" | apt-key add && \
#apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367 && \ # does the same as the previous line
apt-get update && \
apt-get install -y ansible

## add nginx config file
COPY ./nginx_config_site /etc/nginx/sites-available/site
## enable nginx site and remove default site
RUN rm /etc/nginx/sites-enabled/default && ln -s /etc/nginx/sites-available/site /etc/nginx/sites-enabled/
EXPOSE 80/tcp


### add custom files
COPY ./start.sh /start.sh
RUN chmod +x /start.sh


### run custom configuration
RUN chmod o+w /etc/ansible/hosts


### startup command
CMD /start.sh && service nginx start && service php7.0-fpm start