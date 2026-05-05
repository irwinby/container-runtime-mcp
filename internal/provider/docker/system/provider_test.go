package system

import "context"

func nopTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return ctx, func() {}
}
