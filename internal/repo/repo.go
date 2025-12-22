package repo

import "github.com/ManoloEsS/burrow/internal/domain"

type RequestRepo interface {
	Save(r *domain.Request) error
	List() ([]*domain.Request, error)
}

