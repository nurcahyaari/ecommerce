package context

import (
	"context"
	"strconv"
)

type ContextCustomEnum int

const (
	userIdCtx ContextCustomEnum = iota + 1
)

func SetUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userIdCtx, userId)
}

func GetUserId(ctx context.Context) string {
	val := ctx.Value(userIdCtx)
	userId, ok := val.(string)
	if !ok {
		return ""
	}

	return userId
}

func GetUserIdInt64(ctx context.Context) (int64, error) {
	userIdStr := GetUserId(ctx)
	return strconv.ParseInt(userIdStr, 10, 64)
}
