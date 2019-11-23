package webhook_process

type Facade struct {
	webPush ProcessInterface
}

func NewFacade(credFile string) *Facade {
	return &Facade{
		webPush: NewWebPush(credFile),
	}
}

func (f *Facade) ProcessEvent(data string) error {
	err := f.webPush.ProcessEvent(data)
	if err != nil {
		return err
	}
	return nil
}
