package usecase

import (
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/repository"
)

// GuardUseCase is the application service for the file system guard feature.
type GuardUseCase interface {
	Start() error
	Stop()
	Events() <-chan guard.GuardEvent
}

type guardUseCase struct {
	repo repository.GuardRepository
}

// NewGuardUseCase wires the guard use case to its repository implementation.
func NewGuardUseCase(repo repository.GuardRepository) GuardUseCase {
	return &guardUseCase{repo: repo}
}

func (u *guardUseCase) Start() error {
	return u.repo.Start()
}

func (u *guardUseCase) Stop() {
	u.repo.Stop()
}

func (u *guardUseCase) Events() <-chan guard.GuardEvent {
	return u.repo.Events()
}
