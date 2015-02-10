#GoAir
**GoAir** is a multi-platform (Linux/Windows) CLI tool written as a Golang package that implements the [VMware vCloud Air bindings](https://github.com/vmware/govcloudair) package.

It is in early stages right now, so please follow the repo if interested!  The focus up to this point has been in implementing the correct CLI and configuration framework along with basic vCloud Air On-Demand authentication.

The package leverages Steve Francia's [Cobra](https://github.com/spf13/cobra) as a CLI framework and [Viper](https://github.com/spf13/viper) for configuration management.

##Configuration
The project currently has the three areas of configuration possible in order to priority.  

###Flags (--flag)

      --username='username@domain'
      --password='pasword'
      --endpoint="https://us-california-1-3.vchs.vmware.com/"

###Environment Variables (VCLOUDAIR_)

      VCLOUDAIR_USERNAME='username@domain' \
      VCLOUDAIR_PASSWORD='password' \
      VCLOUDAIR_ENDPOINT="https://us-california-1-3.vchs.vmware.com/" \ VCLOUDAIR_SHOW_RESPONSE='true' \
      ./goair

####Configurations files (config.yaml in ~HOME/.goair/ or /etc/goair)

      USERNAME: username@domain
      PASSWORD: password
      ENDPOINT: https://us-california-1-3.vchs.vmware.com/
      SHOW_RESPONSE: false

##Compiling
```git clone https://github.com/emccode/goair```
```go build -i -a```
```go install github.com/emccode/goair```

##Docker Build
Start by cloning the Github repo with ```git clone https://github.com/emccode/goair```.

If you need to build the container manually to get the latest binary and place it will place ```goair``` the binary locally.  This binary may not be executable since it is being compiled for a Linux platform.
```docker run -ti -e REPO_PATH=github.com/emccode/goair -v $(pwd):/output emccode/golang_build_from_url```

Once this is done you can build a new Docker container with this binary.
```docker build -t emccode/goair .```.


##Running
The CLI can be ran as follows.

```goair --help```

You can also leverage the Docker container to run the CLI commands directly or interactively.  To run them directly use the following command.
```docker run -ti -e VCLOUDAIR_USERNAME='username@domain' -e VCLOUDAIR_PASSWORD='password' emccode/goair --help```

An interactive run would look like this ```docker run -ti -e VCLOUDAIR_USERNAME='username@domain' -e VCLOUDAIR_PASSWORD='password' --entrypoint=/bin/bash emccode/goair``` followed by ```goair```.


##CLI Hierarchy
This will be filled out as there are more things added.

      goair ondemand plans get


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
        | us-california-1-3.vchs.vmware.com | 41400e74-4445-49ef-90a4-98da4ccfb16c | Virtual Private Cloud OnDemand | com.vmware.vchs.compute |
        | us-virginia-1-4.vchs.vmware.com   | feda2919-32cb-4efd-a4e5-c5953733df33 | Virtual Private Cloud OnDemand | com.vmware.vchs.compute |
        | uk-slough-1-6.vchs.vmware.com     | 62155213-e5fc-448d-a46a-770c57c5dd31 | Virtual Private Cloud OnDemand | com.vmware.vchs.compute |
        | PMP                               | e79a3b5f-92bc-4b33-aa6c-9dbd1e9d9bfe | My Subscriptions               | com.vmware.vchs.vcim    |
        +-----------------------------------+--------------------------------------+--------------------------------+-------------------------+


Licensing
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

Support
-------
Please file bugs and issues at the Github issues page. For more general discussions you can contact the EMC Code team at <a href="https://groups.google.com/forum/#!forum/emccode-users">Google Groups</a> or tagged with **EMC** on <a href="https://stackoverflow.com">Stackoverflow.com</a>. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.
