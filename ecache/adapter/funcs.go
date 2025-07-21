package adapter

import "context"

type Func = func(ctx context.Context) (value interface{}, err error)
