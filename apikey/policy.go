package apikey

import "github.com/breathbath/goalert/permission"

// GQLPolicy is a GraphQL API key policy.
type GQLPolicy struct {
	Version int
	Query   string
	Role    permission.Role
}
