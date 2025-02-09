package common

import "context"

type LsnID struct{}

func ContextWithID(parent context.Context, value int64) context.Context {
	return context.WithValue(parent, LsnID{}, value)
}

func GetIDFromContext(ctx context.Context) int64 {
	return ctx.Value(LsnID{}).(int64)
}
