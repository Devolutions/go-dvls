package dvls

import (
	"context"
	"fmt"
)

const (
	userListEndpoint        = "/api/security/users/list"
	applicationListEndpoint = "/api/security/application/users/list"
)

var ErrUserNotFound = fmt.Errorf("user not found")
var ErrMultipleUsersFound = fmt.Errorf("multiple users found")

type Users service

// User represents a DVLS user or application account. Application accounts
// have AuthenticationType UserAuthenticationApplication and their Name is the
// application key.
type User struct {
	Id                 string                 `json:"id"`
	Name               string                 `json:"name"`
	FullName           string                 `json:"fullName"`
	Email              string                 `json:"email"`
	AuthenticationType UserAuthenticationType `json:"authenticationType"`
	IsAdministrator    bool                   `json:"isAdministrator"`
	IsEnabled          bool                   `json:"isEnabled"`
	UserGroups         []string               `json:"userGroups"`
}

func userName(u User) string { return u.Name }

// List returns all users, excluding application accounts.
func (c *Users) List() ([]User, error) {
	return c.ListWithContext(context.Background())
}

// ListWithContext returns all users, excluding application accounts.
// The provided context can be used to cancel the request.
func (c *Users) ListWithContext(ctx context.Context) ([]User, error) {
	return fetchDataList[User](ctx, c.client, userListEndpoint, "users")
}

// ListApplications returns all application accounts.
func (c *Users) ListApplications() ([]User, error) {
	return c.ListApplicationsWithContext(context.Background())
}

// ListApplicationsWithContext returns all application accounts.
// The provided context can be used to cancel the request.
func (c *Users) ListApplicationsWithContext(ctx context.Context) ([]User, error) {
	return fetchDataList[User](ctx, c.client, applicationListEndpoint, "applications")
}

// GetByName returns a single user based on its login name.
// Returns ErrUserNotFound if no user is found.
// Returns ErrMultipleUsersFound if more than one user matches the name.
func (c *Users) GetByName(name string) (User, error) {
	return c.GetByNameWithContext(context.Background(), name)
}

// GetByNameWithContext returns a single user based on its login name.
// Returns ErrUserNotFound if no user is found.
// Returns ErrMultipleUsersFound if more than one user matches the name.
// The provided context can be used to cancel the request.
func (c *Users) GetByNameWithContext(ctx context.Context, name string) (User, error) {
	users, err := c.ListWithContext(ctx)
	if err != nil {
		return User{}, err
	}

	return singleByName(users, name, userName, ErrUserNotFound, ErrMultipleUsersFound)
}

// GetApplicationByName returns a single application account based on its application key.
// Returns ErrUserNotFound if no application is found.
// Returns ErrMultipleUsersFound if more than one application matches the name.
func (c *Users) GetApplicationByName(name string) (User, error) {
	return c.GetApplicationByNameWithContext(context.Background(), name)
}

// GetApplicationByNameWithContext returns a single application account based on its application key.
// Returns ErrUserNotFound if no application is found.
// Returns ErrMultipleUsersFound if more than one application matches the name.
// The provided context can be used to cancel the request.
func (c *Users) GetApplicationByNameWithContext(ctx context.Context, name string) (User, error) {
	apps, err := c.ListApplicationsWithContext(ctx)
	if err != nil {
		return User{}, err
	}

	return singleByName(apps, name, userName, ErrUserNotFound, ErrMultipleUsersFound)
}
