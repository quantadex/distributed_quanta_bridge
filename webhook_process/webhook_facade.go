package webhook_process

type Facade struct {
	webPush ProcessInterface
}

func NewFacade() *Facade {
	return &Facade{
		webPush: NewWebPush(),
	}
}

func (f *Facade) ProcessEvent(data string) error {
	err := f.webPush.ProcessEvent(data)
	if err != nil {
		return err
	}
	return nil
}
