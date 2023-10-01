FROM debian:bookworm

LABEL maintainer="info@enforge.de"


## update os & add necessary apt repositories
RUN apt-get update && \
apt-get upgrade -y && \
apt-get autoremove -y && \
apt-get clean -y && \
rm -rf /tmp/* /var/tmp/* /var/cache/apt/* /var/cache/distfiles/*


## install required software
### genisoimage replaces / contains mkisofs
RUN apt-get install -y git binutils gcc liblzma-dev make mtools genisoimage perl syslinux


## get ipxe src
RUN git clone git://git.ipxe.org/ipxe.git


## configure ipxe compile settings
### informations taken from https://github.com/SuperClutch/docker-ipxe-build/blob/master/Dockerfile
RUN sed -i "s%//#define POWEROFF_CMD%#define POWEROFF_CMD%" /ipxe/src/config/general.h


## startup command
CMD cd ipxe/src; \
  make bin-x86_64-efi/ipxe.efi; \
  cp bin-x86_64-efi/ipxe.efi /build/; \
  make bin/undionly.kpxe; \
  cp bin/undionly.kpxe /build/;