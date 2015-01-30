# mkclouddrive
This portion of *Core2F* allows you to leverage a Docker container or Golang binary to create a Cloud Drive ISO and respective Cloud Config file.

The script is based on Kelsey Hightower's Kubeconfig project [https://github.com/kelseyhightower/kubeconfig](Here).

It is in early stages, so it is meant to only allow for SSH public key configuration and a specific output file.

## Build Docker Container
1. Clone the repo with ```git clone https://github.com/emccode/core2f```.
2. Enter the repo and script directory with ```cd core2f/mkclouddrive/```.
3. ```docker build -t emccode/mkclouddrive .```.

## Run Docker Container
There are two volumes that can be mounted.  The first volume specifies where the SSH public key resides, and the second is for the output ISO file.
```docker run -v /Users/clintonkitson/.ssh:/ssh -v $(pwd):/output -ti emccode/mkclouddrive -outfile="/output/configdrive.iso" -publickeyfile="/ssh/id_rsa.pub```
