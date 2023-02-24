# go-dvls
[![Go Reference](https://pkg.go.dev/badge/github.com/Devolutions/go-dvls.svg)](https://pkg.go.dev/github.com/Devolutions/go-dvls)
[![testing](https://github.com/Devolutions/go-dvls/actions/workflows/test.yml/badge.svg)](https://github.com/Devolutions/go-dvls/actions/workflows/test.yml)

:warning: **This client is a work in progress, expect breaking changes between releases** :warning:

Go client for DVLS

Heavily based on the information found on the [Devolutions.Server](https://github.com/Devolutions/devolutions-server/tree/main/Powershell%20Module/Devolutions.Server) powershell module.

## Usage
- Run go get `go get github.com/Devolutions/go-dvls`
- Add the import `import "github.com/Devolutions/go-dvls"`
- Setup the client (we recommend using an [Application ID](https://helpserver.devolutions.net/webinterface_applications.html?q=application+id))
``` go
package main

import (
	"log"

	"github.com/Devolutions/go-dvls"
)

func main() {
    // We strongly recommend using an Application ID with your client
	c, err := dvls.NewClient("username", "password", "https://your-dvls-instance.com")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(c.ClientUser.Username)
}
```

## Documentation
All our documentation is available on [![Go Reference](https://pkg.go.dev/badge/github.com/Devolutions/go-dvls.svg)](https://pkg.go.dev/github.com/Devolutions/go-dvls)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
