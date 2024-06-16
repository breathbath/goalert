package apikey

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/breathbath/goalert/permission"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Middleware struct{}

var _ graphql.OperationParameterMutator = Middleware{}

func (Middleware) ExtensionName() string                          { return "GQLAPIKeyMiddleware" }
func (Middleware) Validate(schema graphql.ExecutableSchema) error { return nil }

func (Middleware) MutateOperationParameters(ctx context.Context, rc *graphql.RawParams) *gqlerror.Error {
	p := PolicyFromContext(ctx)
	if p == nil {
		return nil
	}

	if rc.Query == "" {
		// Allow query to be omitted for API key requests,
		// since they are always fixed to the key itself.
		//
		// The stored query hass been validated beforehand against
		// the API key's embedded hash.
		//
		// This helps with things like key rotations, where the
		// query may not be known to the client, or would otherwise
		// be difficult to update.
		rc.Query = p.Query
	}

	if p.Query != rc.Query {
		return &gqlerror.Error{
			Err:     permission.Unauthorized(),
			Message: "wrong query for API key",
			Extensions: map[string]interface{}{
				"code": "invalid_query",
			},
		}
	}

	return nil
}
