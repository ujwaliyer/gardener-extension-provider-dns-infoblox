# _Downloading Infoblox DDI on VMware vSphere_

This is a step-by-step guide on how to setup the Infoblox DDI software on a vSphere VM. 

The pre-requisite for this process is having a NIOS template file (.OVA file) downloaded on your machine.

## _Instance Setup:_
1.	On the vSphere client go to: Inventories->Hosts and Clusters->Select the resource pool you want to install the NIOS instance on.
2.	Right click the resource pool and select the “Deploy OVF Template” option in the dropdown. 
3.	Follow through the deployment wizard and rename the instance as per your requirement.
4.	When you reach the configuration tab select the model of the NIOS appliance as per your requirement, you can check the various
specs of the different models [here](https://docs.infoblox.com/space/NVIG/35786250#AboutInfobloxNIOSVirtualApplianceforVMware-SupportedvNIOSforVMwareApplianceModels). 
For example, as per the current requirement, a model that is supported as a grid master and grid master candidate with at least 
250GB of disk space was required and thus IB-V815 was selected. The image below shows the licenses that will be available once 
the instance is setup, so choosing one of the below models is required if a temporary license needs to be generated.


![Available Licenses](https://docs.infoblox.com/rest/api/content/35483668/child/attachment/att51356609/download?download=true "Available licenses")

5.	On the select networks tab, select the network you want to assign to the NIOS appliance in the destination column.
6.	You can leave the customize template as is as we will be updating the IP address (LAN1) directly from the NIOS CLI later.
7.	Once you verify the details of configuration via the Ready to Complete section, you can click Finish and this will complete 
the setting up of the instance.


## _Grid Setup:_
1.	Launch the web console from vSphere for the instance that has been set up.
2.	Login to the Infoblox system, the default login credentials will be:

    •	User: admin

    •	Password: Infoblox

3.	Once logged in, a Infoblox> prompt will appear and now you can enter the command set temp_license.
4.	Select 2 for DNS, DHCP, and Grid Licenses and enter y for all the following prompts.
5.	Enter the command set temp_license again and select 4 for the NIOS license, choose the version chosen before in step 4 of 
Instance setup. Note: The model selected here needs to be the same as the model for which you have provisioned the NIOS 
virtual instance in step 4 of instance setup otherwise performance issues might occur.
6.	To confirm that the licenses have been generated, you can enter the show license command to see whether the licenses for DNS, 
DHCP, Grid, and the specific NIOS model are there as seen below

![Generated Licenses](https://github.com/shankmankumar/testing/blob/main/showlicenses.png "generated licenses")

7.	Now we need to setup the LAN1 IP address, for this we need to first find an IP address that is still free within the IP 
address pools set up in NSX-T. For our current purposes we used/will be using 10.16.198.17 as the SAP internal public IP address and 
10.250.0.10 as the internal vSphere IP address. (An explanation of how this was done is given after the setup steps)
8.	These two IP addresses need to be mapped to each other in the NAT section of NSX-T. 
9.	To do this, go to Networking->Networking Services->NAT.
10.	Choose the gateway (gardener-dev-shoot--core--vs-tst1) and then click the Add NAT Rule button.
11.	Enter any name of your choice, for the Action choose DNAT. [You can go [here](https://medium.com/networks-security/nat-snat-dnat-pat-port-forwarding-b7982fab02cd) for more information on the types of NAT]
12.	Leave the source IP address field blank and enter the SAP internal public IP address that was chosen (10.16.198.17 for us) in the Destination IP field.
13.	Enter the vSphere internal IP address (10.250.0.10 for us) in the Translated IP field. Then click on Save. 
14.	Your mapping should look as seen below (ui-v1-nios): 

![NAT mapping](https://github.com/shankmankumar/testing/blob/main/natmappingquality.png "NAT mapping")

15.	Now that the mapping has been done, you can go back to the Infoblox CLI and enter the command set network.
16.	Here, for the IPv4 Address, enter the internal vSphere IP address you want to use (10.250.0.10). For the rest of the entries, 
you can skip and let the default values be populated such that your network settings look as seen below when you use the command 
show network.

![show network](https://github.com/shankmankumar/testing/blob/main/shownetwork.png?raw=true "show network")

17.	Now, when you enter 10.16.198.17 into your browser window, since it has been mapped to 10.250.0.10, which is the IP address of the vSphere VM, you will be able to see the Infoblox Grid Manager directly via your browser.
18.	Upon first opening the Infoblox Grid Manager you will have to accept the terms and then set up the grid master via the Grid 
Setup Wizard. If for some reason, an issue occurs and the grid wizard that automatically shows up when you first open the UI goes 
away, you can find the Grid Setup Wizard as follows: From the Grid tab->Grid Manager->Members tab then expand the toolbar on the 
right and click Grid Properties->Setup (Grid Setup Wizard). 
19.	After following through the Setup Wizard and setting up the Grid Master, your Infoblox DDI installation is complete, and 
the licenses generated can be used for a 60 day period as is. 

## Method Used for Mapping IP Addresses

To determine which IP address was free to be used as the IP for the Infoblox Grid Manager, a previously existing mapping was used. 
The translated IP address from the previous Infoblox-DDI setup (10.250.0.9) was mapped to a range of IP addresses from the range 
10.16.198.* and those IP addresses were entered into a browser till the UI from 10.250.0.9's NIOS instance was seen on the 
browser. Once this is done, we know that that particular IP address (in our case 10.16.198.17) is free to be used and is not currently 
mapped to, or being used by, anything else. We then gave the translated IP as the next IP in the range 10.250.0.*, which in our 
case was 10.250.0.10 and assigned the same to the NIOS instance on vSphere and thus were able to successfully map the internal 
vSphere IP to the SAP Internal Public IP. 