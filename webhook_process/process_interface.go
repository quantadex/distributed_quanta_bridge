package webhook_process

type ProcessInterface interface {
	ProcessEvent(event string) error
}
