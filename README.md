# CORE2F
This project contains relevant steps, automation tools instructions, and eventually Docker containers to help minimize friction when deploying CoreOS to vSphere/vCloud Air/vCloud Director.  


- [Overview](#overview)
- [Manual Method - vCloud Air/vCloud Director](#manual_method)
 - [Networking](#networking) (Optional)
 - [Create cloud-config File and Config Drive ISO](#create_cc_cd)
 - [Upload ISOs to vCloud Director](#upload_isos)
 - [Deploy Template or VApp](#deploy_template)
 - [Attach Media](#attach_media)
 - [Boot VApp](#boot_vapp)
 - [Stop VApp](#stop_vapp) (Optional)
 - [Check-In VApp to Catalog](#checkin_vapp) (Optional)
 - [Deploy VApp](#deploy_vapp) (Optional)
 - [SSH to CoreOS](#ssh) (Optional)
- [Automation](#Automation)
 -[end-to-end](#End-to-End Deploys)
 -[further_customization](#Further Customization)
#<a id="overview">Overview</a>

Have you ever competed for the quickest time to solve a rubix cube?  ```cube2f``` is a special cube that has rounded edges and core lubrication.  The special cube minimizes friction to help a player get an edge.  The combination of this strategy for minimizing friciton for deploying CoreOS and the special rubix cube are where the project name, ```core2f``` comes from.

In quickly changing infrastructure world, ```CoreOS``` is one of the first Linux distributions focused on ditching the legacy baggage that have held common Linux distros back.  The new OS is called a ```Container OS```, where it is hyper-focused on core ingredients to running containers that are currently managed by Docker.  This includes a continuos deployment flow where updates to Stable, Beta, and Alpha channels are hapenning on an ongoing basis.  The Containers themselves are what minimize friction, but getting them deployed in all the necessary places with the latest version may be challening at times.

The most popular Enterprise based Hypervisor is VMware's vSphere.  This hypervisor also serves at the heart of VMware's Public Cloud and On-Demand service as well as vCloud Director for Private Clouds.  Since the Private and vCloud Director sides are mostly considered self-managed clouds it is important to keep up to date CoreOS images available in the catalog.

The project is meant to allow continuous deployment of CoreOS templates or images to VMware infrastructure.  



# <a id="manual_method">Manual Method - vCloud Air/vCloud Director</a>
This section describes the manual method that you can use to establish CoreOS images in your vCloud Air or vCloud Director catalog.  There are plenty of methods out there, but so far they have likely leveraged insecure methods with default certificates to bring CoreOS images up.  This may be fine in some case, but probably isn't a good idea for production desires.  VMware is in tech-preview with CoreOS as a secure version as an OVA which includes the open source vmtools version.
> Warning: the details here may be overkill for your use case since CoreOS can be very easy to get up and running in many cases, ie. boot_kernel coreos.autologin to skip authentication from the console, or from the insecure image with a pre-created public key.  The following process is intended to model a end-to-end automated process that delivers updated CoreOS images and customized cloud-config files through vCloud Director.

The desire for these steps is not only to deploy CoreOS with a custom certificate, but also configure network services in vCloud Air that will expose the VMs similar to other Public Clouds.  For this we will include a Organizational Network/Edge Gateway configuration example from vCloud Air.

One more note.  Certain portions of the steps may be useful on a recurring basis once you have things to a certain point.  For example, once you have created an ISO with your specific public key, you can check that image into the catalog and do deploys from that template without the need to reference this information again.

## <a id="networking">1. (Optional) Networking</a>
The steps listed here are meant to mimic CoreOS deployments in common Public Clouds.  This means the ```Config Drive``` based customization, DHCP/manually configured IPs, and also availability from a publicly accessible IP.  Under the covers, there is a good amount of automation to achieve this.

For the networking example, I am using the new vCloud Air On-Demand service.  The default behavior of the service establishes an Edge Gateway (firewall from the internet to your VDC) and an Organizational Network (connection to your edge and default L2 network).

In order to expose the CoreOS image that you deployed to the internet for usage, or further customization you likely must establish three things:
1. DHCP for the Organizational Network.
2. NAT translation inbound and outbound from the edge.
3. Firewall rules that allow specific or all traffic from the edge.

### (Optional) Create VDC
For this walk through, I am going to create a VDC from scratch.

1. Login and open a chosen vCloud Air VPC
2. Press the ```+``` button next to ```Virtual Data Centers```.
3. Enter any name and continue.  On the left side you should see the new VDC, wait until it is complete (a few minutes).  You might need to refresh your screen.

## Configure Gateway and Network
Once the VDC is created, you must configure the ```Gateway```/```Edge Gateway``` and the ```Network```/```Organizational Network```.

First we will create a ```Public IP``` and inbound/outbound NAT rules.  This ensure traffic can pass from the internet to the VM and from the VM to the internet.  This whole section is optional depending on your requirements.

### Add Public IP
1. Select the VDC
2. Press ```Gateways```.
3. Press the ```GATEWAY ON VDC_NAME```.
4. This should bring up the ```NAT Rules``` Tab.  
5. Press the ```Public IPs``` tab and select ```Add IP Address``` and ```Add```.  You should see ```Adding public IP to gateway``` task at the top.  Once it is done continue.
6. Select the ```Public IPs``` tab again to refresh and highlight the new IP and copy it to your clipboard.

### Add NAT Rules
1. Press ```Add NAT Rule```.
2. Select ```DNAT``` to create the internal to external rule.  This could be created only once for the VDC since all outbound traffic could be sharing a single IP.
3. Select the dropdown for ```Original (External) IP``` and select the new public IP.
4. Leave the reset default and enter ```192.168.109.0/24``` in the ```Translated (Internal) IP/Range dialog```.
5. Press ```Next```.
6. Press ```Add``` to add an additional rule.
7. Select ```SNAT```.
8. Enter ```192.168.109.2``` as the IP.  This can be anything, but we will use the ```192.168.109.2-250``` range later for DHCP.
9.  For ```Translated (External) Source
10. Press ```Next``` followed by ```Finish```.

You have created a dynamic NAT rule that ensures all newly initiated outbound traffic leaves via a specific pulic IP. You have also established that the IP will forward all new traffic to this public IP inbound to ```192.168.109.2```.  The next step is the firewall rules.

### Add Firewall Rules
1. Select the ```Firewall Rules``` tab.  
2. Select ```Add Firewall Rule```.
3. Enter ```ssh``` in the name.
4. Select ```TCP``` for ```Protocol```.
5. Select ```Any``` for ```Source```.
6. Select ```Any``` for ```Source Port```.
7. Select ```External``` for ```Destination```.
8. Enter ```22``` for ```Destination Port```.
9. Press  ```Next```.
10. (Optional) Press ```Add```.
11. Enter ```inside_out``` as the Name.
12. Select ```Source``` as ```Internal```.

The NAT section can be modified if you wish to use a single IP and forward ports such as ```1001 -> 22 at VM1``` and ```1002 -> 22 at VM2```.  This is not suggested with CoreOS since there is already likely dynamic NATing occuring from Docker and will add a layer of complication.  The best method would be to create a new public IP for each VM.

### Remove Static IPs
1. In the navigation bar that the top go to ```Gateways```.
2. Press the ```Manage in vCloud Director``` button.
3. Press the ```Org VDC Networks``` tab.
4. Right click the ```default-routed-network``` and go to ```Properties```.
5. Select ```Network Specification```.
6. Enter ```8.8.8.8``` into ```Primary DNS``` and ```8.8.4.4``` into ```Secondary DNS```.
7. Select the ```192.168.109.2-192.168.109.253``` IP range and remove it or modify it.  If using DHCP this will not be relevant, but if not set properly will cause problems when you enable the same subnet for DHCP.
8. Press ```OK```.

### Enable DHCP
DHCP is needed unless you specify a static IP in your ```Cloud Drive```.  vCloud Director cannot customize the IP for you currently.

1. Select the ```Edge Gateways tab```.
2. Right click ```gateway``` and select ```Edge Gateway Services...```.
3. On the ```DHCP``` tab, press the ```Enable DHCP``` checkbox.
4. Select ```Add```.
5. Ensure ```Enable pool``` is selected and the ```default-routed-network``` is selected.
6. Enter ```192.168.109.2-192.168.109.250``` in the ```IP Range```.
7. Press ```OK```.




## <a id="create_cc_cd">Create cloud-config File and Config Drive ISO</a>
vCloud Director generally makes use of "Guest Customization" that allows you to modify specific parameters in VMs and run custom scripts.  This process is heavily based on vm tools, and isn't very flexibly.  The generally accepted public cloud method is to use cloud-config files that contain YAML definitions of parameters, units, and other details for the guest.  The cloud-config file is traditionally applied via a URL or ISO.  Since currently, the guest customization is ignored for the open source VM tools, this example leverages an ISO that gets mounted to the VM, and when the VM boots it applies the customization which should include the updated certificate and/or ```core``` password.

### Create configdrive.iso
1. From OS X - Install CDR tools package which includes ```ovftool```.  If you have ```brew``` installed, then leverage the following command ```brew install cdrtools```.  There are plenty of other methods to get the ovftool binary installed in most Linux OSs.
2. Create a directory structure to be encapsulated in an ISO with ```mkdir -p new-drive/openstack/latest```
3. Obtain an existing pre-shared key, possibly ```cat ~/.ssh/id_rsa.pub```.
4. Edit user_data file ```vi new-drive/openstack/latest/user_data```.  Change the ```hostname``` or remove the hostname line (if making a template) and modify the ```ssh_authorized_keys```.
```
#cloud-config
hostname: hostname
ssh_authorized_keys:
    - ssh-rsa contents_of_your_public_key_file
```


5. (optional) If you want to set the password instead you can run ```openssl password -1``` and enter a password.  The output can be added at the bottom of the ```user_data``` file.
```
users:
   – name: core
     passwd: key_from_openssl_command
```


6. (optional) Add anything you would like to the ```cloud-config``` file, see the CoreOS page https://coreos.com/docs/cluster-management/setup/cloudinit-config-drive/.
7. Run ```mkisofs -R -V config-2 -o configdrive.iso new-drive/``` to create the ```configdrive.iso``` or respective name that you will later upload to vCloud Director.

## <a id="upload_isos">3. Upload ISOs to vCloud Director</a>
You may only need to do this once if you keep your ```user_data``` file with minimal information.  Once your CoreOS image is online and reachable via SSH, you may at any time re-run ```sudo coreos-cloudinit -from-file user_data``` to add additional details.

### From the GUI - Config Drive
1. (Optional) Login to vCloud Air.
2. Open vCloud Director whether from a ```Virutal Private Cloud``` (Networks -> Manage in vCloud Director) in vCloud Air, or directly from an ```Organization``` in vCloud Director.
3. Press the ```Catalogs``` tab.
4. Press ```Media & Other```.
5. Press the ```Upload...``` icon.
6. Select ```Local file``` and press ```Browse```.  
7. Select the proper ```ISO file```.
8. Press ```Upload```.

### From CLI - Config Drive
1. Install ```ovftool```
 1. (Optional) Download and Install ```ovftool``` from VMware.  A simple google search should yield a proper version.  Right now ```4.0``` is current.
 2. (Optional) Leverage the Docker container process from https://github.com/emccode/docker-ovftool.  This may be handy if you build the Docker container and reuse specific versions or reuse the containers from Public Cloud VMs to get better throughput.
2. Upload the ```Config Drive``` ISO as media to the vCloud Director catalog.
> "/Applications/VMware Fusion.app/Contents/Library/VMware OVF Tool/ovftool"  --sourceType="ISO" --vCloudTemplate="false" configdrive.iso  "vcloud://youraccount@yourdomain@fqdn_region_vca:443?org=your_org_id&vdc=your_vdc_name&media=cyour_iso_name&catalog=default-catalog".

### Repeat the desired process for CoreOS ISO
From the GUI you would either upload the ISO as a VApp or a VApp Template into a catalog.  From the CLI you have both options as well.

### From CLI - CoreOS ISO
1. (Optional) Use the following command to import the OVF as a VApp Template.  Notice the ```alpha``` lists the Alpha channel, which could be Beta, or Stable once they have the same OVA.  The ```current``` references the most recent in that channel, but could be any version with the OVA.
> "/Applications/VMware Fusion.app/Contents/Library/VMware OVF Tool/ovftool"  --acceptAllEulas http://alpha.release.core-os.net/amd64-usr/current/coreos_production_vmware_ova.ova "vcloud://youraccount@yourdomain@fqdn_region_vca:443?org=your_org_id&vdc=your_vdc_name&catalog=default-catalog&vappTemplate=coreos_template_name"
2. (Optional) The following command will import the OVF as a VApp.
> "/Applications/VMware Fusion.app/Contents/Library/VMware OVF Tool/ovftool"  --acceptAllEulas http://alpha.release.core-os.net/amd64-usr/current/coreos_production_vmware_ova.ova "vcloud://youraccount@yourdomain@fqdn_region_vca:443?org=your_org_id&vdc=your_vdc_name&vapp=coreos_vapp_name"

## <a id="deploy_template">4. (Optional) Deploy VApp Template or VApp</a>
If you uploaded as a VApp previously then you can skip this step.  From the vCloud Air or vCloud Director GUI's create a VApp from the template that you uploaded.  Attach the template to the ```default-routed-network```.

## <a id="attach_media">5. Attach Media</a>
In order to customize the CoreOS image with appropriate login information or other details you must attach the ISO that was uploaded.  This ISO can be reused across many CoreOS images, a minimal ISO is desired.  From the vCloud Director GUI, navigate to the VApp.  Right click it, and press ```Open```.  This should reveal the ```Virtual Machine```.  Right click VM and press ```Insert CD/DVD from Catalog...```.  Choose the appropriate media.

## <a id="boot_vapp">6. Boot VApp</a>
Power on the VApp, ensure the VM inside of the VApp is attached to the ```default-routed-network```.  You should see the VM customized with the appropriate hostname if you open the console.  

If the desire is to simply create a good template then you can skip to the ```Add to Template``` section since once the CoreOS VM has booted it will have modified the SSH key.  
### 6. (Optional) Is DHCP configured already?
You can posssibly take it from here.  DHCP is the default behavior for CoreOS unless modified from the ```Cloud Drive``` file to be static.  You can see the IP from vCloud Director listed at the VM, or from the console of the VM at the user prompt.  If it is not there, then you probably need to check out the networking section further since the VM was likely unable to get a DHCP address.

## <a id="stop_vapp">7. (Optional) Stop VApp</a>
Go to the VApps and right click the running VApp and press ```Stop VApp```.  This will initiate a guest shutdown.  Next ejecti the media by opening the VApp and right clicking the VM and selectin ```Eject media```.

## <a id="checkin_vapp">8. (Optional) Check In VApp to Catalog</a>
Right click the VApp and press ```Add to Catalog```.  Fill in desired options, and ensure ```Customize VM Settings``` is selected.  This will make sure that hardware is customized for the VM.

## <a id="deploy">9. (Optional) Deploy VApp</a>
At this point you can deploy your VApp from the template.


## <a id="ssh">10. (Optional) SSH to the VM</a>
In order to access the VM you can now SSH to it.  You had the option of exposing the VM publicly through that NAT and firewall process listed above, or to leave it interal to your vCD networking.  You also had the option to leverage DHCP or static IP addressing.  In addition you could have also specified a password for the ```core``` user or decided to leverage the ```public/private key pair``` method to access the VM without a username.  One last option was to leverage NAT to forward specific ports in bound or expose a single public IP to a  single CoreOS instance.

### SSH with username and specified password
```ssh core@ip```

### SSH with public/private key pair
Depending on your keypair, the file will vary here.  But based on our example, the following would be valid.
```ssh core@external_ip -i ~/.ssh/id_rsa```

### SSH with custom port forward
The ```22222``` port would be changeable of course to whichever port you decided to forward in.
```ssh core@external_ip -i ~/.ssh/id_rsa -p 22222```.

# <a id="automation">Automation</a>

## <a id="end-to-end">End-to-end deploys</a>
To be continued.. Expect Vagrant box examples once the vCloud Air and Director plugins are updated..

## <a id="further_customization">Further customization</a>
It is a good idea to minimize the configuration data in your ```user_data``` file that gets added to the ```Cloud Drive ISO```.  This will make your template more useable.

Once a deploy occurs, and the CoreOS instance has networking access and is accessible via SSH, further customization can take place via SSH.  A good suggestion is to continue to leverage the ```cloud-config``` method.  In this case you can run ```/usr/bin/coreos-cloudinit --from-file /usr/share/oem/cloud-config.yml``` or the same from a URL to further customize the guest after deployment. 


Licensing
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

Support
-------
Please file bugs and issues at the Github issues page. For more general discussions you can contact the EMC Code team at <a href="https://groups.google.com/forum/#!forum/emccode-users">Google Groups</a>. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.
