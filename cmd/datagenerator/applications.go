package main

import (
	"github.com/lab5e/lospan/pkg/keys"
	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/storage"
)

func generateApplications(count int, datastore *storage.Storage, keyGen *keys.KeyGenerator, callback func(generatedApp model.Application)) {
	for i := 0; i < count; i++ {
		app := model.NewApplication()
		var err error
		app.AppEUI, err = keyGen.NewAppEUI()
		if err != nil {
			lg.Error("Unable to generate app EUI. Using random EUI")
			app.AppEUI = randomEUI()
		}
		if err := datastore.CreateApplication(app); err != nil {
			lg.Error("Unable to store application: %v", err)
		} else {
			callback(app)
		}
	}
}
