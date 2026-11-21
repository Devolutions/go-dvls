package dvls

import (
	"context"
	"fmt"
	"net/url"
)

const (
	userEndpoint                = "user"
	applicationIdentityEndpoint = "application-identity"
)

var ErrUserNotFound = fmt.Errorf("user not found")
var ErrMultipleUsersFound = fmt.Errorf("multiple users found")

type Users service

// IdentityProvider is the provider that authenticates a user.
type IdentityProvider string

const (
	IdentityProviderBuiltin         IdentityProvider = "Builtin"
	IdentityProviderContractor      IdentityProvider = "Contractor"
	IdentityProviderActiveDirectory IdentityProvider = "ActiveDirectory"
	IdentityProviderEntraId         IdentityProvider = "EntraId"
	IdentityProviderOkta            IdentityProvider = "Okta"
	IdentityProviderPingOne         IdentityProvider = "PingOne"
	IdentityProviderKeycloak        IdentityProvider = "Keycloak"
)

// User represents a DVLS user. Name is the username. IsWorkspaceOwner reflects the user's own account only, not
// ownership held through a user group.
type User struct {
	Id               string           `json:"id"`
	Name             string           `json:"username"`
	FullName         string           `json:"fullName"`
	Email            string           `json:"email"`
	IdentityProvider IdentityProvider `json:"identityProvider"`
	IsWorkspaceOwner bool             `json:"isWorkspaceOwner"`
	IsEnabled        bool             `json:"isEnabled"`
}

// ApplicationIdentity represents a DVLS application identity. AppKey is the unique key it logs in with; Name is its
// display name, which several identities can share.
type ApplicationIdentity struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	AppKey           string `json:"appKey"`
	IsWorkspaceOwner bool   `json:"isWorkspaceOwner"`
	IsEnabled        bool   `json:"isEnabled"`
}

// List returns all users. Application identities are not included.
func (c *Users) List() ([]User, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all users. Application identities are not included.
// The provided context can be used to cancel the request.
func (c *Users) ListWithContext(ctx context.Context) ([]User, error) {
	return c.list(ctx, nil)
}

func (c *Users) list(ctx context.Context, query url.Values) ([]User, error) {
	users, err := listAll[User](ctx, c.client, query, userEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error while fetching users: %w", err)
	}

	return users, nil
}

// Get returns a single user based on userId.
func (c *Users) Get(userId string) (User, error) {
	return c.GetWithContext(context.Background(), userId)
}

// GetWithContext returns a single user based on userId.
// The provided context can be used to cancel the request.
func (c *Users) GetWithContext(ctx context.Context, userId string) (User, error) {
	reqUrl, err := c.client.v2Url(nil, userEndpoint, userId)
	if err != nil {
		return User{}, fmt.Errorf("failed to build user url: %w", err)
	}

	var user User
	err = c.client.getJSON(ctx, reqUrl, &user)
	if err != nil {
		return User{}, fmt.Errorf("error while fetching user: %w", err)
	}

	return user, nil
}

// GetByName returns a single user based on its username. The name matches ignoring case.
// Returns ErrUserNotFound if no user is found.
// Returns ErrMultipleUsersFound if more than one user matches the name.
func (c *Users) GetByName(name string) (User, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single user based on its username. The name matches ignoring case.
// Returns ErrUserNotFound if no user is found.
// Returns ErrMultipleUsersFound if more than one user matches the name.
// The provided context can be used to cancel the request.
func (c *Users) GetByNameWithContext(ctx context.Context, name string) (User, error) {
	users, err := c.list(ctx, url.Values{"username": {name}})
	if err != nil {
		return User{}, err
	}

	return singleByName(users, name, func(u User) string { return u.Name }, ErrUserNotFound, ErrMultipleUsersFound)
}

// ListApplications returns all application identities.
func (c *Users) ListApplications() ([]ApplicationIdentity, error) {
	return c.ListApplicationsWithContext(context.Background())
}

// ListApplicationsWithContext returns all application identities.
// The provided context can be used to cancel the request.
func (c *Users) ListApplicationsWithContext(ctx context.Context) ([]ApplicationIdentity, error) {
	return c.listApplications(ctx, nil)
}

func (c *Users) listApplications(ctx context.Context, query url.Values) ([]ApplicationIdentity, error) {
	apps, err := listAll[ApplicationIdentity](ctx, c.client, query, applicationIdentityEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error while fetching application identities: %w", err)
	}

	return apps, nil
}

// GetApplication returns a single application identity based on applicationId.
func (c *Users) GetApplication(applicationId string) (ApplicationIdentity, error) {
	return c.GetApplicationWithContext(context.Background(), applicationId)
}

// GetApplicationWithContext returns a single application identity based on applicationId.
// The provided context can be used to cancel the request.
func (c *Users) GetApplicationWithContext(ctx context.Context, applicationId string) (ApplicationIdentity, error) {
	reqUrl, err := c.client.v2Url(nil, applicationIdentityEndpoint, applicationId)
	if err != nil {
		return ApplicationIdentity{}, fmt.Errorf("failed to build application identity url: %w", err)
	}

	var app ApplicationIdentity
	err = c.client.getJSON(ctx, reqUrl, &app)
	if err != nil {
		return ApplicationIdentity{}, fmt.Errorf("error while fetching application identity: %w", err)
	}

	return app, nil
}

// GetApplicationByName returns a single application identity based on its application key (AppKey), which is
// unique; the display name (Name) is not. The key matches ignoring case.
// Returns ErrUserNotFound if no application identity is found.
func (c *Users) GetApplicationByName(appKey string) (ApplicationIdentity, error) {
	return c.GetApplicationByNameWithContext(context.Background(), appKey)
}

// GetApplicationByNameWithContext returns a single application identity based on its application key (AppKey),
// which is unique; the display name (Name) is not. The key matches ignoring case.
// Returns ErrUserNotFound if no application identity is found.
// The provided context can be used to cancel the request.
func (c *Users) GetApplicationByNameWithContext(ctx context.Context, appKey string) (ApplicationIdentity, error) {
	apps, err := c.listApplications(ctx, url.Values{"appKey": {appKey}})
	if err != nil {
		return ApplicationIdentity{}, err
	}

	return singleByName(apps, appKey, func(a ApplicationIdentity) string { return a.AppKey }, ErrUserNotFound, ErrMultipleUsersFound)
}
