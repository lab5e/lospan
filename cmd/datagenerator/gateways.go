package main

import (
	"math"
	"math/rand"

	"github.com/lab5e/lospan/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/storage"
)

func generateGateways(count int, datastore *storage.Storage) []model.Gateway {
	var gws []model.Gateway
	for i := 0; i < count; i++ {
		newGW := model.NewGateway()
		newGW.Latitude = float32(rand.Intn(180)) / math.Pi
		newGW.Longitude = float32(rand.Intn(360)) / math.Pi
		newGW.Altitude = 1.0
		newGW.GatewayEUI = randomEUI()
		newGW.IP = randomIP()
		newGW.StrictIP = rand.Int()%2 == 0
		if err := datastore.CreateGateway(newGW); err != nil {
			lg.Error("Unable to store gateway: %v", err)
		} else {
			gws = append(gws, newGW)
		}
	}
	return gws
}
