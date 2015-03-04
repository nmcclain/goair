#GoAir
**GoAir** is a multi-platform (OS X/Docker/Linux/Windows/FreeBSD) CLI tool.  It's current focus is to simplify the process of deploying machines and configuring network services for on-demand vCloud Air compute services.

See Youtube videos here.



- [Runtime Choices](#runtimechoices)
- [Configuration](#configuration)
  - [Flags with possible runtime persistence](#flagspersistency)
  - [Configuration files](#configfiles)
  - [Environment Variables](#env)
- [Running - Basic](#running-basic)
- [Running - Docker Mini (scratch)](#advanced-docker)
  - [Basic Docker](#basic-docker)
  - [Advanced Docker](#advanced-docker)
- [Deploy VApp Steps](#deployvapp)
- [CLI Command Examples](#cliexamples)
  - [ondemand](#ondemand)
  - [orgvdcnetwork](#orgvdcnetwork)
  - [compute](#compute)
  - [catalog](#catalog)
  - [edgegateway](#edgegateway)
  - [vapp](#vapp)
  - [media](#media)
- [Examples](#examples)
- [About the CLI](#aboutthecli)
- [Future](#future)
- [Contributions](#contributions)
- [Licensing](#licensing)
- [Support](#support)




## <a id="runtimechoices">Runtime Choices</a>
There are plenty of options to install and run *Goair*.  Choose from any one of the following options.  
- Release Binaries - https://github.com/emccode/goair/releases/tag/v0.1.150301
- Ubuntu Goair Docker container (206MB)
- Mini (scratch) Goair Docker container (9MB)
- Clone the github repo and build yourself
  - git clone https://github.com/emccode/goair
    - go build -i -a (to be goair binary for local platform)
  - OR
    - go get github.com/emccode/goair
    - docker run --rm -it -v $GOPATH:/go -w /go/src/github.com/emccode/goair golang:1.4.2-cross make release
    - binaries available in the release/ dir


## <a id="configuration">Configuration</a>
The project currently has the four areas of configuration possible, listed in order of priority.  The combination of flags (parameters), environment variables, and configuration files allows for about any use case possible.

###<a id="flagspersistency">Flags with possible runtime persistence</a>
Flags are simple the ```--``` followed by a paramter.  Certain flags like username and password may only be needed during an initial llogin since there is an element of persistence through Go binary files that save authentication tokens acorss CLI executions.

      --username='username@domain'
      --password='pasword'
      --endpoint="https://us-california-1-3.vchs.vmware.com/api"

###<a id="configfiles">Configurations files</a> (config.yaml in ~HOME/.goair/ or /etc/goair)
The home directory is translated depending on the operating system.  For OS X/Linux the ```HOME``` environment variable is used.  For Windows the ```HOMEDRIVE/HOMEPATH``` combination of environment variables are used unless it is blank, otherwise ```USERPROFILE``` is used.  The next option is the ```/etc/goair``` directory which works across operating systems.

      insecure: 'false'
      username: username@domain
      password: password
      endpoint: https://us-california-1-3.vchs.vmware.com/api


###<a id="env">Environment Variables</a>
Set the environment variables to ```true``` on the boolean based variables.


      VCLOUDAIR_USERNAME: Your vCloud Air username
      VCLOUDAIR_PASSWORD: Your vCloud Air password
      VCLOUDAIR_ENDPOINT: Your preferred vCloud Air intiial endpoint, ie. https://us-california-1-3.vchs.vmware.com/api
      VCLOUDAIR_SHOW_RESPONSE: Whether to show the response or not from the API call
      VCLOUDAIR_SHOW_BODY: Whether to show the HTTP body.  Will intercept POST bodies.
      VCLOUDAIR_INSECURE: Whether to disregard SSL errors
      VCLOUDAIR_SHOW_FLAG: Display details regarding the configuration via Viper for flags, env varibles, configuration files, and gob.
      VCLOUDAIR_SHOW_GOB: Show decoding and encoding details for Go binary files
      CLUE_DEBUG: Show file locations during gob operations.




## <a id="running-basic">Running - Basic</id>
The CLI can be ran as follows.  Using the proper binary from the ```release``` directory the following command will work.

```goair --help```

You can also leverage the Docker container to run the CLI commands directly or interactively.  To run them directly use the following command.
```docker run -ti -e VCLOUDAIR_USERNAME='username@domain' -e VCLOUDAIR_PASSWORD='password' emccode/goair --help```


## <a id="basic-docker">Running - Docker</a>
A great option for running *goair* is through a Docker container.  There are a couple of choices for this.  If you would like to have an interactive session with goair inside of a Docker container you can use the standard ```emccode/goair``` Docker image.  This would be executed as ```docker run -ti emccode/goair```.  From there all of the methods, ie. flags, environment variables, and configuration files are available.

Docker containers can also take advantage of a couple of things.  You can specify ahead of time the environment variables to be used via ```-e VCLOUDAIR_USERNAME=test@test.com``` flags or even in a custom Docker image with ```ENV VCLOUDAIR_USERNAME xxyz```.  This makes the interactive usage of the CLI easier.  In addition you can also map a local directory with the ```config.yaml``` file with a ```-v /Users/username/.goair/:/etc/goair``` flag (or respective to your system).

## <a id="advanced-docker">Running - Docker Mini (Scratch)</a>
The ```goair-mini``` image is a minimal Docker container based on the scratch image.  This means the only space consumed by the container is the *goair* binary file.  The upside to this is the minimal method for distribution.  The downside is that it means there is no interactive usage inside of a container since there is no ```bash```.  You can leverage this style, but you must do as specified prior to get proper configuration to goair as well as mount a temp directory so the go binary files can persist across containers.  You can map these to whichever location you want with ```-v /tmp/:/tmp```.

This method is interactive as well, but from outside the container.  This means you continually execute the ```docker run``` command.  Specifying a ```--rm``` as a flag will ensure the container gets deleted when after command completion.

##<a id="deployvapp">Deploy VApp Steps</a>
The following steps are mostly operational but useful to see a complete flow of getting a VApp deployed from a catalog and operational.

### Login and Choose Compute Resources (VDC)

    goair ondemand login
    goair ondemand plans get | grep region
    goair ondemand login use compute --region= --vdcname=VDC4

### Get VDC Network Name

    goair orgvdcnetwork get

### Get and Deploy Catalog Item

    goair catalog get
    goair catalog get --catalogname="Public Catalog"
    goair catalog get --catalogname="Public Catalog" --catalogitemname=CENTO
    goair catalog deploy --catalogname="Public Catalog" --catalogitemname=CENTO --vmname=centos01

### Get IP
    goair vapp get --vappname=centos01

### Get Available Public IPs and Add NAT

    goair edgegateway get publicip
    goair edgegateway new-natrule 1to1 --externalip=107.189.92.154 --internalip=192.168.109.2 --description=newrule

### Add Firewall Rules for Inbound and Outbound

    goair edgegateway new-firewallrule --destinationport="22" --sourceport="Any" --destinationip="107.189.92.154" --sourceip="Any" --protocol=tcp --description="outside_in"
    goair edgegateway new-firewallrule --destinationport="Any" --sourceport="Any" --destinationip="Any" --sourceip="192.168.109.0/24" --protocol=tcp --description="inside_out"

### Update VM Size and Customization Script

    goair vapp update --vappname=test8 --memorysizemb=2048 --cpucount=4
    goair vapp update guestcustomization script --vappname=test8 < guestCustomizationExample.sh

### Poweron

    goair vapp action poweron --vappname=vappname

### Get Initial Password

    goair vapp get guestcustomization --vappname=vappname

### SSH



##<a id="cliexamples">CLI Command Examples</a>
This will be filled out as there are more things added.

### <a id="ondemand">ondemand</a>

      goair ondemand login
      goair ondemand plans get
      goair ondemand serivcegroupids get
      goair ondemand instances get
      goair ondemand users get
      goair ondemand billable costs get --servicegroupid=4fde19a4-7621-428e-b190-dd4db2e158cd

### <a id="orgvdcnetwork">orgvdcnetwork</a>

      goair orgvdcnetwork get
      goair orgvdcnetwork get --networkname=default-routed-network

### <a id="compute">compute</a>
      goair compute get
      goair compute get --region=us-california-1-3.vchs.vmware.com
      goair compute use --planid=41400e74-4445-49ef-90a4-98da4ccfb16c
      goair compute use --region=us-california-1-3.vchs.vmware.com --name=VDC4

### <a id="catalog">catalog</a>
      goair catalog get
      goair catalog get --catalogname="Public Catalog"
      goair catalog get --catalogname="Public Catalog" --catalogitemname="CentOS64-64Bit"
      goair catalog get vapptemplate --catalogname="Public Catalog" --catalogitemname="CentOS64-64Bit"
      goair catalog deploy --catalogname="Public Catalog" --catalogitemname="CentOS64-64Bit" --vmname="Test2" --vdcnetworkname=default-routed-network
      goair catalog deploy --catalogname="Public Catalog" --catalogitemname="CentOS64-64Bit" --vmname="Test2" --vdcnetworkname=default-routed-network --runasync=true

### <a id="edgegateway">edgegateway</a>
      goair edgegateway get

      goair edgegateway new-natrule 1to1 --externalip=107.189.92.154 --internalip=192.168.109.2 --description=newrule
      goair edgegateway remove-natrule 1to1 --externalip=107.189.92.154 --internalip=192.168.109.2
      goair edgegateway get natrule
      goair edgegateway get gatewayinteface
      goair edgegateway get iprange
      goair edgegateway get publicip
      goair edgegateway new-publicip --publicipcount=3 --networkname=d3p4v54-ext
      goair edgegateway remove-publicip --networkname=d3p4v54-ext --publicip=107.189.87.208
      goair edgegateway new-firewallrule --destinationport="22" --sourceport="Any" --destinationip="107.189.92.154" --sourceip="Any" --protocol=tcp --description="outside_in"

      goair edgegateway new-firewallrule --destinationport="Any" --sourceport="Any" --destinationip="Any" --sourceip="192.168.109.0/24" --protocol=tcp --description="inside_out"
      goair edgegateway new-firewallrule --destinationport="Any" --sourceport="Any" --destinationip="107.189.92.154" --sourceip="Any" --protocol=icmp --description="outside_in_icmp"
      goair edgegateway remove-firewallrule --ruleid=1

### <a id="vapp">vapp</a>

      goair vapp get
      goair vapp get --vappname=test8
      goair vapp get-status --vappname=test8
      goair vapp get --vappid=urn:vcloud:vapp:789d295e-296f-4679-94a4-c17ba36c3d62
      goair vapp get vm --vappname=test8
      goair vapp update --vappname=test8 --memorysizemb=2048 --cpucount=4
      goair vapp get guestcustomization --vappname=test8
      goair vapp update guestcustomization script --vappname=test8 < guestCustomizationExample.sh
      goair vapp action poweron --vappname=vappname
      goair vapp action poweroff --vappname=vappname
      goair vapp action reboot --vappname=vappname
      goair vapp action reset --vappname=vappname
      goair vapp action suspend --vappname=vappname
      goair vapp action shutdown --vappname=vappname
      goair vapp action undeploy --vappname=vappname
      goair vapp action deploy --vappname=vappname
      goair vapp action delete --vappname=vappname
      goair vapp insertmedia --vappname=vappname --medianame=configdrive-basic-id_rsa6.iso
      goair vapp ejectmedia --vappname=vappname --medianame=configdrive-basic-id_rsa6.iso

### <a id="media">media</a>

      goair media get

##<a id="examples">Examples</a>
Here is the help screen that is available at every level using ```help``` or ```--help```.

    Usage:
      goair [flags]
      goair [command]

    Available Commands:
      ondemand                  ondemand
      compute                   compute
      vapp                      vapp
      catalog                   catalog
      orgvdcnetwork             orgvdcnetwork
      edgegateway               edgegateway
      help [command]            Help about any command

     Available Flags:
          --Config="": config file (default is $HOME/goair/config.yaml)
          --help=false: help for goair


##<a id="output">Output from commands</a>
The intended output from the commands is to be a format that is both human readable and interpretable in a structured programmatic way.  For this, YAML has been chosen for most command outputs.  


    ./goair compute use --region=us-california-1-3.vchs.vmware.com --vdcname=VDC4
    href: https://us-california-1-3.vchs.vmware.com/api/compute/api/vdc/cbecb4b5-4267-4018-9458-a05d56936eff
    id: ""
    type: application/vnd.vmware.vcloud.vdc+xml
    name: VDC4
    rel: down



## <a id="aboutthecli">About the CLI</a>
Since we are using Go, the first major benefit is that we are able to cross-compile it to a bunhc of different platforms and architectures.  The following list covers the binaries that are compiled along with their relative sizes.

      9941784 goair-Darwin-i386*
     12340576 goair-Darwin-x86_64*
     12275632 goair-FreeBSD-amd64*
      9860304 goair-FreeBSD-i386*
      9890552 goair-Linux-armv6l*
      9890552 goair-Linux-armv7l*
      9917272 goair-Linux-i386*
      9917272 goair-Linux-i686*
      8762528 goair-Linux-static*
     12281200 goair-Linux-x86_64*
      9997312 goair.exe*

The *goair* application functions identically across any of these executables.  There are a couple of notable differences.

- SSL Certificates - platforms with ca-certificates in non-default locations or not installed must use the environment variable ```VCLOUDAIR_USECERTS=true```.  This will foce the usage of default ca-certificates.
- Configuration via environment variables or location of configuration files may differ across platforms.  See below.

In order to make the CLI as easy to use as possible you can expect certain things (auth tokens) to be cached in local temp locations.  This makes it possible to run commands like ```goair use compute --region=there --vdcname=vdc1``` and have all further commands respect this context.  This functionality is driven by the [Clue package](https://github.com/emccode/clue).

Finally, the since *Goair* is a compiled binary and possible even static (zero dependencies), it is extremely efficient to use interactively and simple to distribute.

##<a id="future">Future</a>
- Upload to Catalog
- Other vCA services

##<a id="contributions">Contributions</a>
The package leverages a handful of open source technologies and projects.

- [Go Programming Language](http://golang.org/pkg/)
- [VMware vCloud Air Go API bindings](https://github.com/vmware/govcloudair)
- [Cobra](https://github.com/spf13/cobra) CLI framework
- [Viper](https://github.com/spf13/viper) for configuration management.
- [Golang-Crosscompile](https://github.com/davecheney/golang-crosscompile.git)
- [Gotablethis](https://github.com/emccode/gotablethis)
- [Clue](https://github.com/emccode/clue)


<a id="licensing">Licensing</a>
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

<a id="support">Support</a>
-------
Please file bugs and issues at the Github issues page. For more general discussions you can contact the EMC Code team at <a href="https://groups.google.com/forum/#!forum/emccode-users">Google Groups</a> or tagged with **EMC** on <a href="https://stackoverflow.com">Stackoverflow.com</a>. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.
