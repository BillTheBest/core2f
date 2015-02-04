# vcloud_upload_media
This portion of *Core2F* represents functionality that allows you to upload an OVA or OVFs as Virtual Machines and their respective disks directly to a VApp.

There is currently a Ruby script that can act in a standalone manner along with a Docker container that works without dependencies.

The ```vcloud_upload_vapp.rb``` script is based on VMware's vCloud Air ```ruby_vcloud_sdk``` project.

## Build Docker Container
1. Clone the repo with ```git clone https://github.com/emccode/core2f```.
2. Enter the repo and script directory with ```cd core2f/vcloud_upload_vapp/```.
3. ```docker build -t emccode/vcloud_upload_vapp .```.

## Run Docker Container
If your file is in an OVA format you must first untar it to a specific directory so there is an OVF and VMDKs present.  You can do this with ```tar -zxvf file.ova -c outputdir```.

The following command should be self explanatory for the environemnt variables that must be specified.  Ensure that you are specifying a correct path with the ```-v $(pwd):/host``` parameter where ```$(pwd)``` represents the directory the media to be uploaded exists.
> docker run -e vcloud_url='https://us-virginia-1-4.vchs.vmware.com:443' -e vcloud_username='user@domain@org_name' -e vcloud_password='vcloud_password' -e vcloud_vdc_name='vdc_name' -e vcloud_ovf_directory='./coreos_production_vmware_ova/' -e vcloud_vapp_name='test13' -ti emccode/vcloud_upload_vapp 
