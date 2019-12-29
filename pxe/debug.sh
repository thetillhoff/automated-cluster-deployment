#!/bin/bash

docker run --rm --cap-add=NET_ADMIN --net=host -it -v $PWD/containerdata:/container container_pxeansible bash -c "/start.sh && service nginx start && service php7.0-fpm start; echo 'run dnsmasq --no-daemon --log-queries --log-facility=/container/dnsmasq.log'; bash"


# --rm  removes the container after is exits
# -i    keep stdin open even if not attached
# -t    allocate a pseudo-tty
# -a    attach to stdin, stdout and/or stderr
# -v    attaches volume (source:destination)
# --volume      -v
#   use either of the following container-internal paths
#   /container
#   /mnt/... (nothing in it by default)
#   /data
# -h    container host name
# --hostname    -h
# --name        assign name to conatiner
# --cap-add=NET_ADMIN       allow interaction with network stack
# --net=host                add host-networking stack to container -> port 53 of host becomes DNS and is accessible from the rest of the network