package restapi

import (
	"net/http"

	"github.com/lab5e/lospan/pkg/model"
)

func (s *Server) getApplication(w http.ResponseWriter, r *http.Request) *model.Application {
	appEUI, err := euiFromPathParameter(r, "aeui")
	if err != nil {
		http.Error(w, "Malformed Application EUI", http.StatusBadRequest)
		return nil
	}
	application, err := s.context.Storage.GetApplicationByEUI(appEUI)
	if err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return nil
	}
	return &application
}
