# vcloud_upload_media
This portion of *Core2F* represents functionality that allows you to upload ISO media files to a specific catalog and VDC.  Currently ```ovftool``` does not allow this capability unless you create catalogs manually that align to storage policies from specific VDCs.

There is currently a Ruby script that can act in a standalone manner along with a Docker container that works without dependencies.

The ```vcloud_upload_media.rb``` script is based on VMware's vCloud Air ```ruby_vcloud_sdk``` project.

## Build Docker Container
1. Clone the repo with ```git clone https://github.com/emccode/core2f```.
2. Enter the repo and script directory with ```cd core2f/vcloud_upload_media/```.
3. ```docker build -t emccode/vcloud_upload_media```.

## Run Docker Container
The following command should be self explanatory for the environemnt variables that must be specified.  Ensure that you are specifying a correct path with the ```-v $(pwd):/host``` parameter where ```$(pwd)``` represents the directory the media to be uploaded exists.
> docker run -e vcloud_url='https://us-virginia-1-4.vchs.vmware.com:443' -e vcloud_username='user@domain@org_name' -e vcloud_password='password' -e vcloud_catalog_name='default-catalog' -e vcloud_vdc_name='VDC3' -e vcloud_dest_iso_name='media4.iso' -e vcloud_file='/host/configdrive.iso' -v $(pwd):/host -ti emccode/vcloud_upload_media
