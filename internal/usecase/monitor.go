package usecase

type MonitorUseCase interface {
	Start(directory string, recursive bool) error
	Stop() error
}
