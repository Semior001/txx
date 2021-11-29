package txx

import (
	"context"
)

type runnerKey struct{}

func putRunner(ctx context.Context, r Runner) context.Context {
	return context.WithValue(ctx, runnerKey{}, r)
}

func getTx(ctx context.Context) Runner {
	if ctx == nil {
		return nil
	}
	rr := ctx.Value(runnerKey{})
	if rr == nil {
		return nil
	}
	return (rr).(Runner)
}
