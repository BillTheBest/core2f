# CORE2F
Docker containers and relevant automation tools that minimize friction when deploying CoreOS to vSphere/vCloud Air/vCloud Director.  

Overview
--------

Have you ever competed for the quickest time to solve a rubix cube?  ```cube2f``` is a special cube that has rounded edges and core lubrication.  The special cube minimizes friction to help a player get an edge.

In the infrastructure world, ```CoreOS``` is one of the first Linux distributions focused on ditching the legacy baggage that have held common Linux distros back.  The new OS is called a ```Container OS```, where it is hyper-focused on core ingredients to running containers that are currently managed by Docker.  This includes a continuos deployment flow where updates to Stable, Beta, and Alpha channels are hapenning on an ongoing basis.  The Containers themselves are what minimize friction, but getting them deployed in all the necessary places with the latest version may be challening at times.

The most popular Enterpirse based Hypervisor is VMware's vSphere.  This hypervisor also serves at the heart of VMware's Public Cloud and On-Demand service as well as vCloud Director for Private Clouds.  Since the Private and vCloud Director sides are mostly considered self-managed clouds it is important to keep up to date CoreOS images available in the catalog.

The project is meant to allow continuous deployment of CoreOS templates or images to VMware infrastructure.  



Licensing
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

Support
-------
Please file bugs and issues at the Github issues page. For more general discussions you can contact the EMC Code team at <a href="https://groups.google.com/forum/#!forum/emccode-users">Google Groups</a>. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.
