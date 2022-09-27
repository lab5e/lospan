package restapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/lab5e/lospan/pkg/monitoring"

	"github.com/ExploratoryEngineering/logging"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/storage"
)

// Read application from request body. Emits error message to client if there's an error
func (s *Server) readAppFromRequest(w http.ResponseWriter, r *http.Request) (apiApplication, error) {
	buf, err := io.ReadAll(r.Body)
	app := apiApplication{}

	if err != nil {
		logging.Warning("Unable to read request body from %s: %v", r.RemoteAddr, err)
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return app, err
	}

	if err := json.Unmarshal(buf, &app); err != nil {
		logging.Info("Unable to unmarshal JSON: %v, (%s)", err, string(buf))
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return app, err
	}
	return app, nil
}

// Handle POST to application collection, ie create application
func (s *Server) createApplication(w http.ResponseWriter, r *http.Request) {
	application, err := s.readAppFromRequest(w, r)
	if err != nil {
		return
	}

	var overrideEUI bool
	if application.ApplicationEUI != "" {
		overrideEUI = true
		if application.eui, err = protocol.EUIFromString(application.ApplicationEUI); err != nil {
			http.Error(w, "Invalid EUI format", http.StatusBadRequest)
			return
		}
	}
	if !overrideEUI {
		if application.eui, err = s.context.KeyGenerator.NewAppEUI(); err != nil {
			logging.Warning("Unable to generate application EUI: %v", err)
			http.Error(w, "Unable to create application EUI", http.StatusInternalServerError)
			return
		}
	}
	application.ApplicationEUI = application.eui.String()

	// This might seem like a baroque way of getting an EUI but since EUIs can
	// be user-specified we will have EUIs that collide once in a while. Most of
	// the time this shouldn't be an issue but if large blocks are added there
	// might be more than one. This will attempt to create it for a relatively
	// small number of requests and skip the EUI counter forwards 10 steps at a
	// time.
	appErr := storage.ErrAlreadyExists
	app := application.ToModel()

	attempts := 1
	for appErr == storage.ErrAlreadyExists && attempts < 10 {
		appErr = s.context.Storage.CreateApplication(app)
		if appErr == nil {
			break
		}
		if appErr == storage.ErrAlreadyExists && overrideEUI {
			// Can't reuse app EUI
			http.Error(w, "Application EUI is already in use", http.StatusConflict)
			return
		}
		if appErr != storage.ErrAlreadyExists {
			// Some other error - fail with 500
			logging.Warning("Unable to store application: %s", err)
			http.Error(w, "Unable to store application", http.StatusInternalServerError)
			return
		}

		logging.Warning("EUI (%s) for application is already in use. Trying another EUI.", app.AppEUI)
		app.AppEUI, err = s.context.KeyGenerator.NewAppEUI()
		if err != nil {
			http.Error(w, "Unable to generate application identifier", http.StatusInternalServerError)
			return
		}
		attempts++
	}
	if appErr == storage.ErrAlreadyExists {
		logging.Error("Unable to create find available EUI even after 10 attempts. Returning error to client")
		http.Error(w, "Unable to store application", http.StatusInternalServerError)
		return
	}

	// Set the tags property if it isn't already set
	if application.Tags == nil {
		application.Tags = make(map[string]string)
	}

	monitoring.ApplicationCreated.Increment()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(application); err != nil {
		logging.Warning("Unable to marshal application object: %v", err)
	}
}

// Handle GET on application collection, ie list applications
func (s *Server) applicationList(w http.ResponseWriter, r *http.Request) {
	// GET returns a JSON array with applications.
	applications, err := s.context.Storage.ListApplications()
	if err != nil {
		logging.Warning("Unable to read application list: %v", err)
		http.Error(w, "Unable to load applications", http.StatusInternalServerError)
		return
	}
	appList := newApplicationList()
	for application := range applications {
		appList.Applications = append(appList.Applications, newAppFromModel(application))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(appList); err != nil {
		logging.Warning("Unable to marshal application object: %v", err)
	}
}

// Shows a list of applications.
func (s *Server) applicationListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		s.applicationList(w, r)

	case http.MethodPost:
		s.createApplication(w, r)

	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

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

// Return a single application formatted as JSON.
func (s *Server) applicationInfoHandler(w http.ResponseWriter, r *http.Request) {
	application := s.getApplication(w, r)
	if application == nil {
		return
	}

	switch r.Method {

	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newAppFromModel(*application)); err != nil {
			logging.Warning("Unable to marshal application with EUI %s into JSON: %v", application.AppEUI, err)
		}

	case http.MethodDelete:
		err := s.context.Storage.DeleteApplication(application.AppEUI)
		switch err {
		case nil:
			monitoring.ApplicationRemoved.Increment()
			monitoring.RemoveAppCounters(application.AppEUI)
			w.WriteHeader(http.StatusNoContent)
		case storage.ErrNotFound:
			// This is covered above but race conditions might apply here
			http.Error(w, "Application not found", http.StatusNotFound)
		case storage.ErrDeleteConstraint:
			http.Error(w, "Application can't be deleted", http.StatusConflict)
		default:
			logging.Warning("Unable to delete application: %v", err)
			http.Error(w, "Unable to delete application", http.StatusInternalServerError)
		}
		return

	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}
