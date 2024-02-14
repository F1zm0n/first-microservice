package main

import (
	"context"
	"encoding/json"
	"github.com/f1zm0n/broker/event"
	"net/http"
)

func (app *Config) pushToQueue(ctx context.Context, name, msg, severity string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	payload := LogPayload{
		Name: name,
		Data: msg,
	}
	j, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	err = emitter.Push(ctx, string(j), severity)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) logEventRabbit(w http.ResponseWriter, ctx context.Context, l LogPayload) {
	//l.Severity = "log.INFO"
	err := app.pushToQueue(ctx, l.Name, l.Data, "log.INFO") //todo improve that part with severity
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via rabbitMQ"
	app.writeJSON(w, http.StatusAccepted, payload)

}
