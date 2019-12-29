# automated cluster deployment

## what this is about
To pack all required tools for automated deployment of a cluster infrastructure, this repository was created. It includes the generation of ssh keys, a pxe server in dns-proxy mode and an ansible server combined in a container.
It can/will be decided later in the process whether the pxe-server and/or the ansible server will be migrated into the cluster.
The key generation is optional at this point, as by default a public ssh key is loaded from a provided url.

## notes
Most reoccuring tasks are specified in the ```Taskfile.yml```.

## TODO
- [ ] load generated public ssh key to new hosts, either additionally or as the new default
- [ ] setup ansible with new hosts
- [ ] configure ansible for reload config at an specified interval (from a git repo)
- [ ] configure ansible to store new hosts to a git repo