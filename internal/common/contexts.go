package common

import "context"

type lsnID struct {
}

func ContextWithID(parent context.Context, value int64) context.Context {
	return context.WithValue(parent, lsnID{}, value)
}

func GetIDFromContext(ctx context.Context) int64 {
	return ctx.Value(lsnID{}).(int64)
}
