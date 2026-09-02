# Go client for Devolutions Server
[![Go Reference](https://pkg.go.dev/badge/github.com/Devolutions/go-dvls.svg)](https://pkg.go.dev/github.com/Devolutions/go-dvls)
[![testing](https://github.com/Devolutions/go-dvls/actions/workflows/test.yml/badge.svg)](https://github.com/Devolutions/go-dvls/actions/workflows/test.yml)

:warning: **This client is a work in progress, expect breaking changes between releases** :warning:

## Compatibility

| go-dvls version | DVLS version   |
|-----------------|----------------|
| 0.16.0+         | 2026.x         |
| 0.15.0          | 2024.x, 2025.x |

Heavily based on the information found on the [Devolutions.Server](https://github.com/Devolutions/devolutions-server/tree/main/Powershell%20Module/Devolutions.Server) powershell module.

Users, applications and user groups (`client.Users`, `client.UserGroups`) expose the principal ids used as assignees in role assignments and as roles in entry permissions.

## Usage
- Run go get `go get github.com/Devolutions/go-dvls`
- Add the import `import "github.com/Devolutions/go-dvls"`
- Setup the client using either an [Application ID](https://docs.devolutions.net/server/web-interface/administration/security-management/applications/) (App Key + App Secret) or an [API key](https://docs.devolutions.net/server/knowledge-base/how-to-articles/generate-manage-and-use-devolutions-server-api-keys)

Using an Application ID:
``` go
package main

import (
	"log"

	"github.com/Devolutions/go-dvls"
)

func main() {
	c, err := dvls.NewClient("appKey", "appSecret", "https://your-dvls-instance.com")
	if err != nil {
		log.Fatal(err)
	}
	_ = c
}
```

Using an API key:
``` go
package main

import (
	"log"

	"github.com/Devolutions/go-dvls"
)

func main() {
	c, err := dvls.NewClientWithApiKey("apiKey", "https://your-dvls-instance.com")
	if err != nil {
		log.Fatal(err)
	}
	_ = c
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
