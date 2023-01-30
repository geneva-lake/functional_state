package service

import (
	"encoding/json"
	"net/http"

	"github.com/geneva-lake/functional_state/general"
	"github.com/geneva-lake/functional_state/logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//   - -------------------------------------------------------------------------------------------------------------------
//     Make user information endpoint where the chain of collecting user information
//     is started
//   - -------------------------------------------------------------------------------------------------------------------
func MakeUserEndpoint(svc *Service) general.Endpoint {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		defer func() {
			r := recover()
			if r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp := UserResponse{
					Status: general.StatusError,
					Error:  InternalError,
				}
				go logger.Log(logger.Panic, svc.config.Name, http.StatusInternalServerError, nil, r, nil, resp)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}()
		vars := mux.Vars(r)
		varsID, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			resp := UserResponse{
				Status: general.StatusError,
				Error:  NoUserID,
			}
			json.NewEncoder(w).Encode(resp)
			go logger.Log(logger.Error, svc.config.Name, http.StatusBadRequest, nil, nil, nil, resp)
			return
		}
		userID, err := uuid.Parse(varsID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := UserResponse{
				Status: general.StatusError,
				Error:  WrongUserID,
			}
			json.NewEncoder(w).Encode(resp)
			go logger.Log(logger.Error, svc.config.Name, http.StatusBadRequest, err, nil, varsID, resp)
			return
		}

		userFlow := svc.WithContext(r.Context()).WithUserID(userID).FromStorage().
			FromFinancials().Answer()
		resp := UserResponse{}
		switch userFlow.Status {
		case StoredUserNotFound:
			resp.Status = general.StatusError
			resp.Error = UserNotFound
			w.WriteHeader(http.StatusBadRequest)
			go logger.Log(logger.Error, svc.config.Name, http.StatusBadRequest,
				nil, nil, userFlow.ID, resp)
		case StoredUserInternalError:
			resp.Status = general.StatusError
			resp.Error = InternalError
			w.WriteHeader(http.StatusInternalServerError)
			go logger.Log(logger.Error, svc.config.Name, http.StatusInternalServerError,
				userFlow.Err, nil, userFlow.ID, resp)
		case TransactionsNotFound, TransactionsReceived:
			resp.Status = general.StatusOK
			dto := NewUserDTO(userFlow)
			resp.Result = dto
			w.WriteHeader(http.StatusOK)
			go logger.Log(logger.Info, svc.config.Name, http.StatusOK,
				nil, nil, userFlow.ID, resp)
		case TransactionsNotAvailable:
			resp.Status = general.StatusOK
			dto := NewUserDTO(userFlow)
			resp.Result = dto
			w.WriteHeader(http.StatusOK)
			go logger.Log(logger.Error, svc.config.Name, http.StatusOK,
				userFlow.Err, nil, userFlow.ID, resp)
		default:
			resp.Status = general.StatusError
			resp.Error = InternalError
			w.WriteHeader(http.StatusInternalServerError)
			go logger.Log(logger.Error, svc.config.Name, http.StatusInternalServerError,
				userFlow.Err, nil, userFlow.ID, resp)
		}
		json.NewEncoder(w).Encode(resp)
	}
}
