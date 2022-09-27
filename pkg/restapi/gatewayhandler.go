package restapi

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/ExploratoryEngineering/logging"
	"github.com/lab5e/lospan/pkg/monitoring"
	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/storage"
)

func (s *Server) gatewayList(w http.ResponseWriter, r *http.Request) {
	gateways, err := s.context.Storage.GetGatewayList()
	if err != nil {
		logging.Warning("Unable to get list of gateways: %v", err)
		http.Error(w, "Unable to read list of gateways", http.StatusInternalServerError)
		return
	}

	gatewayList := newGatewayList()
	for gateway := range gateways {
		gatewayList.Gateways = append(gatewayList.Gateways, newGatewayFromModel(gateway))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gatewayList); err != nil {
		logging.Warning("Unable to marshal gateway list: %v", err)
	}
}

func (s *Server) createGateway(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Info("Unable to read request body: %v", err)
		http.Error(w, "Unable to read request", http.StatusInternalServerError)
		return
	}

	gateway := apiGateway{}
	if err := json.Unmarshal(buf, &gateway); err != nil {
		logging.Info("Unable to unmarshal JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(gateway.GatewayEUI) == "" {
		http.Error(w, "Missing gateway EUI", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(gateway.IP) == "" {
		http.Error(w, "Missing gateway IP", http.StatusBadRequest)
		return
	}

	if gateway.eui, err = protocol.EUIFromString(gateway.GatewayEUI); err != nil {
		http.Error(w, "Invalid gateway EUI", http.StatusBadRequest)
		return
	}
	if gateway.ipaddr = net.ParseIP(gateway.IP); gateway.ipaddr == nil {
		http.Error(w, "Invalid IP format", http.StatusBadRequest)
		return
	}

	// Sanity check the lat/lon coordinates
	if gateway.Longitude > 360 || gateway.Longitude < -360 ||
		gateway.Latitude < -90 || gateway.Latitude > 90 {
		http.Error(w, "Longitude must be [180.0, 180] Latitude must be and [-90, 90]",
			http.StatusBadRequest)
		return
	}

	modelGw := gateway.ToModel()
	if err = s.context.Storage.CreateGateway(gateway.ToModel()); err != nil {
		if err == storage.ErrAlreadyExists {
			http.Error(w, "A gateway with that EUI alread exists", http.StatusConflict)
			return
		}
		logging.Warning("Unable to store gateway with EUI %s: %v", gateway.GatewayEUI, err)
		http.Error(w, "Unable to store gateway", http.StatusInternalServerError)
		return
	}

	monitoring.GatewayCreated.Increment()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newGatewayFromModel(modelGw)); err != nil {
		logging.Warning("Unable to marshal gateway with EUI %s into JSON: %v", modelGw.GatewayEUI, err)
	}
}

func (s *Server) gatewayListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.gatewayList(w, r)

	case http.MethodPost:
		s.createGateway(w, r)

	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func (s *Server) gatewayInfoHandler(w http.ResponseWriter, r *http.Request) {
	eui, err := euiFromPathParameter(r, "geui")
	if err != nil {
		logging.Debug("Unable to convert EUI from string: %v", err)
		http.Error(w, "Invalid EUI", http.StatusBadRequest)
		return
	}

	modelGateway, err := s.context.Storage.GetGateway(eui)
	if err != nil {
		logging.Info("Unable to look up gateway with EUI %s: %v", eui, err)
		http.Error(w, "Gateway not found", http.StatusNotFound)
		return
	}

	switch r.Method {

	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newGatewayFromModel(modelGateway)); err != nil {
			logging.Warning("Unable to marshal gateway with EUI %s into JSON: %v", modelGateway.GatewayEUI, err)
		}
		return

	case http.MethodPut:
		var values map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		alt, ok := values["altitude"].(float64)
		if ok {
			modelGateway.Altitude = float32(alt)
		}
		lat, ok := values["latitude"].(float64)
		if ok {
			modelGateway.Latitude = float32(lat)
		}
		lon, ok := values["longitude"].(float64)
		if ok {
			modelGateway.Longitude = float32(lon)
		}
		ip, ok := values["ip"].(string)
		if ok {
			if modelGateway.IP = net.ParseIP(ip); modelGateway.IP == nil {
				http.Error(w, "Invalid IP address", http.StatusBadRequest)
				return
			}
		}
		strict, ok := values["strictIP"].(bool)
		if ok {
			modelGateway.StrictIP = strict
		}

		if err := s.context.Storage.UpdateGateway(modelGateway); err != nil {
			logging.Warning("Unable to update gateway: %v", err)
			http.Error(w, "Unable to update gateway", http.StatusInternalServerError)
			return
		}
		monitoring.GatewayUpdated.Increment()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newGatewayFromModel(modelGateway)); err != nil {
			logging.Warning("Unable to marshal gateway with EUI %s into JSON: %v", modelGateway.GatewayEUI, err)
		}

	case http.MethodDelete:
		if err := s.context.Storage.DeleteGateway(eui); err != nil {
			logging.Warning("Unable to delete gateway: %v", err)
			http.Error(w, "Unable to remove gateway", http.StatusInternalServerError)
			return
		}
		monitoring.GatewayRemoved.Increment()
		monitoring.RemoveGatewayCounters(eui)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
