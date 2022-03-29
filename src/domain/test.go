package domain

import "context"

type TestRepository interface {
	Test(ctx context.Context) ([]string, []string)
}
