package service

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/notification"
	"chem-factory/internal/repository"
)

type Manager struct {
	auth         *auth.Manager
	repository   *repository.Manager
	notification *notification.Manager
}

func New(a *auth.Manager, r *repository.Manager, n *notification.Manager) *Manager {
	return &Manager{
		auth:         a,
		repository:   r,
		notification: n,
	}
}
