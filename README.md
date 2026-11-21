# Go client for Devolutions Server
[![Go Reference](https://pkg.go.dev/badge/github.com/Devolutions/go-dvls.svg)](https://pkg.go.dev/github.com/Devolutions/go-dvls)
[![testing](https://github.com/Devolutions/go-dvls/actions/workflows/test.yml/badge.svg)](https://github.com/Devolutions/go-dvls/actions/workflows/test.yml)

:warning: **This client is a work in progress, expect breaking changes between releases** :warning:

## Compatibility

| go-dvls version | DVLS version                                                  |
|-----------------|---------------------------------------------------------------|
| 1.0.0+         | The first DVLS release with the public API v2 (`/api/v2`) and later |
| 0.16.0 - 0.22.x | 2026.x                                                        |
| 0.15.0          | 2024.x, 2025.x                                                |

From 1.0.0, go-dvls calls only the DVLS public API v2 (`/api/v2`): it no longer calls `/api/v1` or the internal routes
of the web interface, and it does not work with a DVLS server without `/api/v2`. Stay on 0.22.x for those servers.

Users, application identities and user groups (`client.Users`, `client.UserGroups`) expose the principal IDs used as
assignees in role assignments and in entry permissions.

## Upgrading to 1.0.0

Sessions: an application identity token expires a fixed time after login (about 7 minutes) and v2 has no refresh
operation. The client no longer checks the session before every request: when a request answers 401, it logs in again
once and retries it. Concurrent requests that get a 401 share one login, and a request waiting for it returns when
its own context is done. API key clients (`NewClientWithApiKey`) are unchanged.

Breaking changes:

- Enums are strings carrying the v2 names: `SecurityRoleOverride`, `SecurityRoleRight`, `AdministrativeRoleScopeType`,
  `AdministrativeRoleAssigneeType`, `AdministrativePermission`, `UserGroupType`, `EntryCertificateDataMode` (`Url`,
  `File`). `SaveResult`, `UserAuthenticationType`, `ServerLoginResult`, `ServerConnectionType` and the
  `AdministrativePermission*` constants are removed; check failures with `IsNotFound` and the sentinel errors.
- Renamed constants: `AdministrativeRoleScopeGlobal` → `AdministrativeRoleScopeWorkspace`,
  `AdministrativeRoleScopeOrganizationalUnit` → `AdministrativeRoleScopeBusinessUnit`,
  `AdministrativeRoleAssigneeApplication` → `AdministrativeRoleAssigneeApplicationIdentity`,
  `SecurityRoleRightConnectionHistory` → `SecurityRoleRightEntryHistory`, `SecurityRoleRightHandbook`,
  `SecurityRoleRightEditHandbook`, `SecurityRoleRightDeleteHandbook` → `...Documentation`,
  `UserGroupTypeOffice365` → `UserGroupTypeEntraId`. Removed: the Gateway and PamProvider scopes,
  `SecurityRoleRightSealed`, and the PAM, gateway, workspace settings and user groups built-in role IDs.
- Every entry is reached through its vault: `Entries.Certificate.Get/GetFileContent/Delete`, `Entries.Host.Get`,
  `Entries.Website.Get` and `Entries.Permissions.Get/Set` take the vault ID before the entry ID.
- Entry permissions: `EntrySecurity` is `{RoleOverride, Permissions}`. `ViewOverride` and `ViewRoles` are removed:
  the View right is set in `Permissions` like any other right. `EntryPermission.Roles` holds user, user group and
  application identity IDs. The zero `EntrySecurity` resets an entry.
- Role assignments: `AdministrativeRoleAssignments.New/List/Get/Delete` replace `AddMember`, `GetMembers`,
  `ListScopeMembers`, `ListByAssignee`, `Me`, `UpdateMembers`, `RemoveMember` and `DeleteScope`.
  `AdministrativeRoleAssignment` is one assignment (`Id`, `RoleId`, `AssigneeId`, `AssigneeType`, `ScopeType`,
  `ScopeId`, names); `AdministrativeRoleAssignmentRequest` and the `AdministrativeRoleAssignmentFilter` fields
  (`RoleId`, `AssigneeId`, `ScopeType`, `ScopeId`) replace `AdministrativeRoleMemberRequest`,
  `AdministrativeRoleMember` and the old filter. Creating an existing assignment returns
  `ErrAdministrativeRoleAssignmentAlreadyExists`; a `List` filter naming a role, assignee or scope that does not
  exist fails with an HTTP 400 (`GetMembers` returned an empty list).
- `AdministrativeRole`: `AssignableScopes` replaces `SupportedScopes`; `RequiredPermissionOnScope`, `IsAssignable`,
  `IsPam`, `IsPrivileged` and `IsUsed` are removed.
- `User`: `IdentityProvider` replaces `AuthenticationType` and `IsWorkspaceOwner` replaces `IsAdministrator`;
  `UserGroups` is removed. `Users.ListApplications` and `GetApplicationByName` return `ApplicationIdentity`;
  `GetApplicationByName` matches the unique application key (`AppKey`), not the display name (`Name`).
  `UserGroup`: `IsWorkspaceOwner` replaces `IsAdministrator`, `IsEnabled` is added.
- `EntryCertificate.Expiration` is returned in UTC and an update with a zero `Expiration` clears it.
  `EntryCertificate.Password` is a `*string`: on `Update`, nil keeps the stored password and a pointer to `""` clears
  it; the certificate `Update` returns has the password sent, nil (unknown) when it kept the stored one. `Update`
  no longer reads the password, so editing the metadata of a certificate does not need the right to view it.
  `EntryHost` and `EntryWebsite` lose `ConnectionType`; `EntryWebsite.ConnectionSubType` is the browser and
  `WebBrowserApplication` is removed.
- `Entries.Host.GetHostDetails` and `Entries.Website.GetWebsiteDetails` return the username and password stored on
  the entry, no longer the resolved credentials: for an entry whose credentials are linked, inherited or come from
  elsewhere they are empty or stale. `HostDetails.CredentialMode` / `WebsiteDetails.CredentialMode` and
  `CredentialEntryId` tell where the entry takes its credentials from.
- `GetServerInfo` (version only) replaces `GetPublicServerInfo` and `GetPrivateServerInfo`; `GetServerTimezones` is
  removed. `Response` holds only the body, `Response.CheckRespSaveResult` and `RequestOptions.RawBody` are removed.
- Behaviour: name lookups (`GetByName`, `GetEntriesOptions.Name`) match ignoring case. A vault created with
  `VaultContentTypeDefault` keeps `Default`. Entry tags come back sorted. On an update, the non-secret data fields
  are always sent, so an empty one is cleared, while an empty secret keeps the stored one. A folder's `Path` is its
  parent folder everywhere (it used to read as its own full path): `GetEntriesOptions.Path` lists a folder under its
  parent, not under itself, and an update sending the folder's own full path is refused.
- The entry routes answer 404 for a vault the caller cannot open, the same answer as for an entry that does not
  exist, so an entry 404 alone does not prove the entry is gone: `Vaults.ContentReadable` tells the two apart.

## Usage
- Run go get `go get github.com/Devolutions/go-dvls`
- Add the import `import "github.com/Devolutions/go-dvls"`
- Setup the client using either an [application identity](https://docs.devolutions.net/server/web-interface/administration/security-management/applications/) (App Key + App Secret) or an [API key](https://docs.devolutions.net/server/knowledge-base/how-to-articles/generate-manage-and-use-devolutions-server-api-keys)

Using an application identity:
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

`NewClient` logs in without a deadline. Use `NewClientWithContext` to cancel the initial login or give it a timeout:
``` go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

c, err := dvls.NewClientWithContext(ctx, "appKey", "appSecret", "https://your-dvls-instance.com")
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
