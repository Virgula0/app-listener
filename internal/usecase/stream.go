package usecase

// StreamUseCase is the application service of a mode whose engine only starts, stops and streams
// events. Consumers (the CLI delivery layer) only depend on this interface.
type StreamUseCase[E any] interface {
	Start() error
	Stop()
	Events() <-chan E
}

type streamUseCase[E any] struct {
	repo StreamUseCase[E]
}

func (u *streamUseCase[E]) Start() error {
	return u.repo.Start()
}

func (u *streamUseCase[E]) Stop() {
	u.repo.Stop()
}

func (u *streamUseCase[E]) Events() <-chan E {
	return u.repo.Events()
}
