package service

import (
	"context"

	//%: {{ if .Examples }}
	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema/skeleton"
	//%: ` | replace "project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
	//%: {{ end }}
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

//%: {{ if .Examples }}
func (s *Service) Dummy(ctx context.Context, req *skeleton.DummyRequest) error {
	return nil
} //%: {{ end }}
