package services

import (
	"context"

	"github.com/prosperitybot/worker/domain"
)

type TestService interface {
	Test(ctx context.Context) ([]string, []string)
}

type DefaultTestService struct {
	repo domain.TestRepository
}

func (s DefaultTestService) Test(ctx context.Context) ([]string, []string) {
	return nil, nil
}

func NewTestService(repo domain.TestRepository) DefaultTestService {
	return DefaultTestService{repo}
}
