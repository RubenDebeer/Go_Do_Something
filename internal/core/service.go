package core

import (
	"github.com/RubenDeBeer/Go_Do_Something/internal/port"
)

// Service contains the business use cases and only depends on ports.

/*
This struct defined a field called repo and that field is a interface, which is a type in golang.
This type now needs to implement the functions defined in the interface.
*/
type Service struct {
	repo port.Repository
}

func NewService(r port.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) AddValue(v string) error {
	if v == "" {
		return nil
	}
	return s.repo.Add(v)
}

func (s *Service) ListValues() ([]string, error) {
	return s.repo.List()
}

func (s *Service) DeleteLast() error {
	return s.repo.DeleteLast()
}
