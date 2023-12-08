package handler

import (	
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/gorilla/mux"

	"github.com/go-rbac/internal/core"
	"github.com/go-rbac/internal/erro"
	
)

var childLogger = log.With().Str("handler", "handler").Logger()

func (h *HttpWorkerAdapter) Health(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Health")

	health := true
	json.NewEncoder(rw).Encode(health)
	return
}

func (h *HttpWorkerAdapter) Live(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Live")

	live := true
	json.NewEncoder(rw).Encode(live)
	return
}

func (h *HttpWorkerAdapter) Header(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("Header")
	
	json.NewEncoder(rw).Encode(req.Header)
	return
}

func (h *HttpWorkerAdapter) AddRole( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("AddRole")

	role := core.RoleData{}
	err := json.NewDecoder(req.Body).Decode(&role)
    if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(erro.ErrUnmarshal.Error())
        return
    }

	res, err := h.workerService.PutRole(req.Context(), role)
	if err != nil {
		switch err {
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) GetRole(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("GetRole")

	vars := mux.Vars(req)
	varID := vars["id"]

	role := core.RoleData{}
	role.Role.Name = varID
	
	res, err := h.workerService.GetRole(req.Context(), role)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			rw.WriteHeader(404)
			json.NewEncoder(rw).Encode(err.Error())
			return
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) AddPolicy( rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("AddPolicy")

	policy := core.PolicyData{}
	err := json.NewDecoder(req.Body).Decode(&policy)
    if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(erro.ErrUnmarshal.Error())
        return
    }

	res, err := h.workerService.PutPolicy(req.Context(), policy)
	if err != nil {
		switch err {
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}

func (h *HttpWorkerAdapter) GetPolicy(rw http.ResponseWriter, req *http.Request) {
	childLogger.Debug().Msg("GetPolicy")

	vars := mux.Vars(req)
	varID := vars["id"]

	policy := core.PolicyData{}
	policy.Policy.Name = varID
	
	res, err := h.workerService.GetPolicy(req.Context(), policy)
	if err != nil {
		switch err {
		case erro.ErrNotFound:
			rw.WriteHeader(404)
			json.NewEncoder(rw).Encode(err.Error())
			return
		default:
			rw.WriteHeader(500)
			json.NewEncoder(rw).Encode(err.Error())
			return
		}
	}

	json.NewEncoder(rw).Encode(res)
	return
}
