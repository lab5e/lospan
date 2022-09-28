package main

import (
	"flag"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/server"
	"github.com/lab5e/lospan/pkg/storage"
)

var params struct {
	ConnectionString string
	UserCount        int
}

var defaultMA protocol.MA

func init() {
	flag.IntVar(&params.UserCount, "users", 1000, "Number of users to generate")
	flag.StringVar(&params.ConnectionString, "connectionstring", "postgres://localhost/congress?sslmode=disable", "PostgreSQL connection string")
	flag.Parse()

	var err error
	defaultMA, err = protocol.NewMA([]byte{0, 9, 9})
	if err != nil {
		panic(err)
	}
}

const appsPerUser = 5
const devicesPerApp = 30
const dataPerDevice = 100
const noncesPerDevice = 30
const gatewaysPerUser = 2

func main() {
	lg.SetLogLevel(lg.InfoLevel)
	lg.Info("This is the data generator tool")
	//datastore := memstore.CreateMemoryStorage(0, 0)
	datastore, err := storage.CreateStorage(params.ConnectionString)
	if err != nil {
		lg.Error("Unable to create datastore: %v", err)
		return
	}

	keygen, err := server.NewEUIKeyGenerator(defaultMA, 0, datastore)
	if err != nil {
		lg.Error("Unable to create key generator: %v", err)
		return
	}

	gateways := generateGateways(gatewaysPerUser, datastore)
	generateApplications(appsPerUser, datastore, &keygen, func(createdApplication model.Application) {
		generateDevices(devicesPerApp, createdApplication, datastore, &keygen, func(createdDevice model.Device) {
			generateDeviceData(createdDevice, dataPerDevice, gateways, datastore)
			generateDownstreamMessage(createdDevice, datastore)
			generateNonces(createdDevice, noncesPerDevice, datastore)
		})
	})

}
