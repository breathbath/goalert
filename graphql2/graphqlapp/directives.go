package graphqlapp

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/breathbath/goalert/expflag"
	"github.com/breathbath/goalert/validation"
)

func Experimental(ctx context.Context, obj interface{}, next graphql.Resolver, flagName string) (res interface{}, err error) {
	if !expflag.ContextHas(ctx, expflag.Flag(flagName)) {
		return nil, validation.NewGenericError("experimental flag not enabled")
	}

	return next(ctx)
}
