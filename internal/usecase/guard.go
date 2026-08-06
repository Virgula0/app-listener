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
	// PopulateInodes fills the inode map with the guarded tree's contents.
	// Must be called before Start so every file is protected from the
	// moment the guard is attached (not just those reached lazily).
	PopulateInodes() error
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

func (u *guardUseCase) PopulateInodes() error {
	return u.repo.PopulateInodes()
}
