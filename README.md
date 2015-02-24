#GoAir
**GoAir** is a multi-platform (Linux/Linux Static/OS X/Windows/FreeBSD) CLI tool written as a Golang package that implements the [VMware vCloud Air bindings](https://github.com/vmware/govcloudair) package.

It is in early stages right now, so please follow the repo if interested!  The focus up to this point has been in implementing the correct CLI and configuration framework along with basic vCloud Air On-Demand authentication.

The package leverages Steve Francia's [Cobra](https://github.com/spf13/cobra) as a CLI framework and [Viper](https://github.com/spf13/viper) for configuration management.

##Configuration
The project currently has the four areas of configuration possible, listed in order of priority.  

###Flags with possible runtime persistence

      --username='username@domain'
      --password='pasword'
      --endpoint="https://us-california-1-3.vchs.vmware.com/api"

###Flags saved by ```use``` statement
A ```use``` command can be used in certain scenarios where the specified flags will be saved for usage later.
      goair compute use --region=us-virginia-1-4.vchs.vmware.com
      goair compute use --region=""

####Configurations files (config.yaml in ~HOME/.goair/ or /etc/goair)

      username: username@domain
      password: password
      endpoint: https://us-california-1-3.vchs.vmware.com/api

###Environment Variables (VCLOUDAIR_)

      VCLOUDAIR_USERNAME='username@domain' \
      VCLOUDAIR_PASSWORD='password' \
      VCLOUDAIR_ENDPOINT="https://us-california-1-3.vchs.vmware.com/api" \
      VCLOUDAIR_SHOW_RESPONSE='true' \
      VCLOUDAIR_SHOW_BODY='true' \
      VCLOUDAIR_INSECURE='true' \
      ./goair


##Compiling
    git clone https://github.com/emccode/goair
    go build -i -a
    go install github.com/emccode/goai

Additionally if you want to cross-compile this you can use the following command to create a release directory under the goair folder with the release binaries.

This will create binaries for OS X/FreeBSD/Linux/Linux Static/Windows.

```docker run --rm -it -v $GOPATH:/go -w /go/src/github.com/emccode/goair golang:1.3-cross make release```



##Docker Build
Run the ```Compiling``` steps first.  From there you can issue a ```docker build -t emccode/goair .``` command to build the new Docker container.  This container will be using the ```scratch``` image which will have a size of around 10MB.


##Running
The CLI can be ran as follows.  Using the proper binary from the ```release``` directory the following command will work.

```goair --help```

You can also leverage the Docker container to run the CLI commands directly or interactively.  To run them directly use the following command.
```docker run -ti -e VCLOUDAIR_USERNAME='username@domain' -e VCLOUDAIR_PASSWORD='password' emccode/goair --help```


##CLI Hierarchy
This will be filled out as there are more things added.

      goair ondemand login
      goair ondemand plans get
      goair ondemand serivcegroupids get
      goair ondemand instances get
      goair ondemand users get
      goair ondemand billable costs get --servicegroupid=4fde19a4-7621-428e-b190-dd4db2e158cd
      goair compute get
      goair compute get --region=us-california-1-3.vchs.vmware.com
      goair compute use --planid=41400e74-4445-49ef-90a4-98da4ccfb16c
      goair compute use --region=us-california-1-3.vchs.vmware.com --name=VDC4
      goair vapp get
      goair vapp get --vappname=test8
      goair vapp get --vappid=urn:vcloud:vapp:789d295e-296f-4679-94a4-c17ba36c3d62
      goair vapp action poweron --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action poweroff --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action reboot --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action reset --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action suspend --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action shutdown --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action undeploy --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action deploy --vappname=vApp_clintonskitson@gmail.com_5
      goair vapp action delete --vappname=vApp_clintonskitson@gmail.com_5

##Examples
Here is the help screen that is available at every level using ```help``` or ```--help```.
    Usage of goair:

    Usage:
      goair [flags]
      goair [command]

    Available Commands:
      ondemand                  ondemand
      help [command]            Help about any command

     Available Flags:
          --Config="": config file (default is $HOME/goair/config.yaml)
      -h, --help=false: help for goair

    Use "goair help [command]" for more information about that command.


The following is an example of the current operation and output in a nicely formatted table.

    dicey1:goair clintonkitson$ bin/goair.darwin ondemand plans get
        +-----------------------------------+--------------------------------------+--------------------------------+-------------------------+
        |              REGION               |                  ID                  |              NAME              |       SERVICENAME       |
        +-----------------------------------+--------------------------------------+--------------------------------+-------------------------+
        | us-california-1-3.vchs.vmware.com | 41400e72-4445-49ef-90a4-98da4ccfb16c | Virtual Private Cloud OnDemand | com.vmware.vchs.compute |
        | us-virginia-1-4.vchs.vmware.com   | feda2913-32cb-4efd-a4e5-c5953733df33 | Virtual Private Cloud OnDemand | com.vmware.vchs.compute |
        | uk-slough-1-6.vchs.vmware.com     | 62155211-e5fc-448d-a46a-770c57c5dd31 | Virtual Private Cloud OnDemand | com.vmware.vchs.compute |
        | PMP                               | e79a3b5g-92bc-4b33-aa6c-9dbd1e9d9bfe | My Subscriptions               | com.vmware.vchs.vcim    |
        +-----------------------------------+--------------------------------------+--------------------------------+-------------------------+


Licensing
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

Support
-------
Please file bugs and issues at the Github issues page. For more general discussions you can contact the EMC Code team at <a href="https://groups.google.com/forum/#!forum/emccode-users">Google Groups</a> or tagged with **EMC** on <a href="https://stackoverflow.com">Stackoverflow.com</a>. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.
