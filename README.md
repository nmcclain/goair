#GoAir
**GoAir** is a multi-platform (Linux/Windows) CLI tool written as a Golang package that implements the [VMware vCloud Air bindings](https://github.com/vmware/govcloudair) package.

It is in early stages right now, so please follow the repo if interested!

The package leverages Steve Francia's [Cobra](https://github.com/spf13/cobra) as a CLI framework and [Viper](https://github.com/spf13/viper) for configuration management.

##Configuration
The project currently has the three areas of configuration possible in order to priority.  
- Flags (--flag)


      --username='username@domain'
      --password='pasword'
      --endpoint="https://us-california-1-3.vchs.vmware.com/"

- Environment Variables (VCLOUDAIR_)


      VCLOUDAIR_USERNAME='username@domain' \
      VCLOUDAIR_PASSWORD='password' \
      VCLOUDAIR_ENDPOINT="https://us-california-1-3.vchs.vmware.com/" \ VCLOUDAIR_SHOW_RESPONSE='true' \
      ./goair

- Configurations files (config.yaml in ~HOME/.goair/ or /etc/goair)


      USERNAME: username@domain
      PASSWORD: password
      ENDPOINT: https://us-california-1-3.vchs.vmware.com/
      SHOW_RESPONSE: false

##Compiling
```go install github.com/emccode/goair```

##Docker Build
```docker build -t emccode/goair .```

##Running
The CLI can be run be run the CLI as follows.

```docker run -ti -e VCLOUDAIR_USERNAME='username@domain' -e VCLOUDAIR_PASSWORD='password' emccode/goair --help```
```goair --help```

##CLI Hierarchy
      goair plans get



Licensing
---------
Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

Support
-------
Please file bugs and issues at the Github issues page. For more general discussions you can contact the EMC Code team at <a href="https://groups.google.com/forum/#!forum/emccode-users">Google Groups</a> or tagged with **EMC** on <a href="https://stackoverflow.com">Stackoverflow.com</a>. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.
