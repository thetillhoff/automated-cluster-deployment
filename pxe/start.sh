#!/bin/bash

## allow nginx to create files in /container/web/mac folder
chmod -R o+w /container/web/mac
## allow nginx to read files from /key folder
chmod -R o+r /key
## allow nginx to write to ansible-hosts file
chmod o+w /etc/ansible/hosts
