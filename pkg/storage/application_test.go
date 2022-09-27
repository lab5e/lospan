package storage

import (
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/model"
)

func testApplicationStorage(
	appStorage *Storage,
	t *testing.T) {

	application := model.Application{
		AppEUI: makeRandomEUI(),
	}

	if err := appStorage.CreateApplication(application); err != nil {
		t.Error("Got error adding application: ", err)
	}

	// Rinse and repeat
	if err := appStorage.CreateApplication(application); err == nil {
		t.Error("Shouldn't be able to add application twice: ", err)
	}

	// Open the application
	existingApp, err := appStorage.GetApplicationByEUI(application.AppEUI)
	if err != nil {
		t.Error("Shouldn't get error when opening an application that is added: ", err)
	}
	if !existingApp.Equals(application) {
		t.Error("The application doesn't match the stored one")
	}
	// Try to open an application that doesn't exist
	if _, err = appStorage.GetApplicationByEUI(makeRandomEUI()); err == nil {
		t.Error("Shouldn't be able to open unknown application")
	}

	// Get list of all applications
	found := 0
	appCh, _ := appStorage.ListApplications()
	select {
	case <-appCh:
		found++
	case <-time.After(time.Millisecond * 100):
		t.Error("Did not get any data on app channel")
	}

	if found == 0 {
		t.Error("Did not get any data on app channel")
	}

	if err := appStorage.DeleteApplication(application.AppEUI); err != nil {
		t.Fatalf("Got error deleting application: %v", err)
	}

	if err := appStorage.DeleteApplication(application.AppEUI); err == nil {
		t.Fatal("Expected error when deleting application but didn't get one")
	}
}
