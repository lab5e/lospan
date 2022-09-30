package storage

import (
	"database/sql"
	"testing"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestApplicationStorage(t *testing.T) {
	assert := require.New(t)

	connectionString := ":memory:"
	db, err := sql.Open(driverName, connectionString)
	assert.Nil(err, "No error opening database")
	assert.Nil(createSchema(db), "Error creating db in %s", connectionString)
	db.Close()

	appStorage, err := CreateStorage(connectionString)
	assert.Nil(err, "Did not expect error: %v", err)
	defer appStorage.Close()

	application := model.Application{
		AppEUI: makeRandomEUI(),
	}

	assert.NoError(appStorage.CreateApplication(application))

	// Rinse and repeat
	assert.Error(appStorage.CreateApplication(application))

	// Open the application
	existingApp, err := appStorage.GetApplicationByEUI(application.AppEUI)
	assert.NoError(err, "Shouldn't get error when opening an application that is added")

	assert.True(existingApp.Equals(application))

	// Try to open an application that doesn't exist
	_, err = appStorage.GetApplicationByEUI(makeRandomEUI())
	assert.Error(err)

	// Get list of all applications
	apps, err := appStorage.ListApplications()
	assert.NoError(err)
	assert.Contains(apps, application, "Returned list contains application")

	assert.NoError(appStorage.DeleteApplication(application.AppEUI))

	assert.Error(appStorage.DeleteApplication(application.AppEUI), "Should get error when applications does not exist")
}
