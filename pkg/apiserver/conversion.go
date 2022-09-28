package apiserver

import (
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/pb/lospan"
)

// Convert model.Application -> lospan.Application

func toAPIApplication(app model.Application) *lospan.Application {
	return &lospan.Application{
		Eui: app.AppEUI.String(),
	}
}
