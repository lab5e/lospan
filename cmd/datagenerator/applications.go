package main

import (
	"github.com/ExploratoryEngineering/logging"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

func generateApplications(count int, datastore *storage.Storage, keyGen *server.KeyGenerator, callback func(generatedApp model.Application)) {
	for i := 0; i < count; i++ {
		app := model.NewApplication()
		var err error
		app.AppEUI, err = keyGen.NewAppEUI()
		if err != nil {
			logging.Error("Unable to generate app EUI. Using random EUI")
			app.AppEUI = randomEUI()
		}
		if err := datastore.CreateApplication(app); err != nil {
			logging.Error("Unable to store application: %v", err)
		} else {
			callback(app)
		}
	}
}
