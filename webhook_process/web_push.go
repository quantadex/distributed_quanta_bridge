package webhook_process

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
)

type WebPush struct {
	app    *firebase.App
	client *messaging.Client
	store  *firestore.Client
	ctx    context.Context
}

type Event struct {
	Name   string
	Quanta string
	TxId   string
}

func NewWebPush() *WebPush {
	w := &WebPush{}
	err := w.initialize()
	if err != nil {
		panic(err)
	}
	return w
}

func (w *WebPush) initialize() error {
	//opt := option.WithCredentialsFile("quantadice-01-firebase-adminsdk-35ckj-7577f12803.json")
	ctx := context.Background()
	w.ctx = ctx
	conf := &firebase.Config{DatabaseURL: "https://quantadice-01.firebaseio.com"}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}
	w.app = app

	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}
	w.client = client

	store, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	w.store = store
	return nil
}

func (w *WebPush) ProcessEvent(data string) error {
	var event *Event
	json.Unmarshal([]byte(data), &event)
	dsnap, err := w.store.Collection("notification_tokens").Doc(event.Quanta).Get(w.ctx)
	if err != nil {
		return err
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m["token"])

	registrationToken := m["token"]

	message := &messaging.Message{
		Webpush: &messaging.WebpushConfig{
			Headers: map[string]string{"Authorization": "key=AAAAEvjxLIs:APA91bGGRr4TFFxFVdfKsnVtj3LN1rf4ugY7qTIcz2lZD5a84nr5IaRagcVY9bvTjVh0eeA-P5uyhGcoh38G6X_FGOFL1mp7fGBOkx6xotiJxd3Y09NURTtMI0bU5uWsIqjOJekcTvKv"},
			Notification: &messaging.WebpushNotification{
				Title:      "Transaction Notification",
				Body:       event.Name,
				CustomData: map[string]interface{}{
					//"click_action": url,
				},
			},
		},
		Token: registrationToken.(string),
	}

	resp, err := w.client.Send(w.ctx, message)
	if err != nil {
		return err
	}
	fmt.Println("Successfully sent message:", resp)
	return nil
}
