# automated cluster deployment

## what this is about
To pack all required tools for automated deployment of a cluster infrastructure, this repository was created. It includes the generation of ssh keys, a pxe server in dns-proxy mode and an ansible server combined in a container.
It can/will be decided later in the process whether the pxe-server and/or the ansible server will be migrated into the cluster.
The key generation is optional at this point, as by default a public ssh key is loaded from a provided url.

## notes
Most reoccuring tasks are specified in the ```./Taskfile.yml```.
If further information is needed, please look at ```./pxe/readme.md``` or ```./ansible/readme.md``` respectively.

### timings
To get a feeling for when a problem occured, here are some measured times for a bare-metal-deployment of a single machine:
From pressing the boot button until accessible: ~15min
Poweroff: ~5s
Boot: ~ 40s
