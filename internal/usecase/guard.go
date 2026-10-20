package usecase

import (
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/repository"
)

// GuardUseCase is the application service for the file system guard feature.
type GuardUseCase = StreamUseCase[guard.GuardEvent]

// NewGuardUseCase wires the guard use case to its repository implementation.
func NewGuardUseCase(repo repository.GuardRepository) GuardUseCase {
	return &streamUseCase[guard.GuardEvent]{repo: repo}
}
