package core

import (
	"github.com/RubenDeBeer/Go_Do_Something/internal/port"
)

// Service contains the business use cases and only depends on ports.
type Service struct {
	repo port.Repository
}

func NewService(r port.Repository) *Service {
	return &Service{repo: r}
}
