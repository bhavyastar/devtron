package terminal

import (
	"encoding/json"
	"errors"
	"github.com/devtron-labs/devtron/api/restHandler/common"
	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/devtron-labs/devtron/pkg/clusterTerminalAccess"
	"github.com/devtron-labs/devtron/pkg/user"
	"github.com/devtron-labs/devtron/pkg/user/casbin"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type UserTerminalAccessRestHandler interface {
	StartTerminalSession(w http.ResponseWriter, r *http.Request)
	UpdateTerminalSession(w http.ResponseWriter, r *http.Request)
	UpdateTerminalShellSession(w http.ResponseWriter, r *http.Request)
	FetchTerminalStatus(w http.ResponseWriter, r *http.Request)
	StopTerminalSession(w http.ResponseWriter, r *http.Request)
	DisconnectTerminalSession(w http.ResponseWriter, r *http.Request)
	DisconnectAllTerminalSessionAndRetry(w http.ResponseWriter, r *http.Request)
	FetchTerminalPodEvents(w http.ResponseWriter, r *http.Request)
	FetchTerminalPodManifest(w http.ResponseWriter, r *http.Request)
}

type UserTerminalAccessRestHandlerImpl struct {
	Logger                    *zap.SugaredLogger
	UserTerminalAccessService clusterTerminalAccess.UserTerminalAccessService
	Enforcer                  casbin.Enforcer
	UserService               user.UserService
	validator                 *validator.Validate
}

func NewUserTerminalAccessRestHandlerImpl(logger *zap.SugaredLogger, userTerminalAccessService clusterTerminalAccess.UserTerminalAccessService, Enforcer casbin.Enforcer,
	UserService user.UserService, validator *validator.Validate) *UserTerminalAccessRestHandlerImpl {
	return &UserTerminalAccessRestHandlerImpl{
		Logger:                    logger,
		UserTerminalAccessService: userTerminalAccessService,
		Enforcer:                  Enforcer,
		UserService:               UserService,
		validator:                 validator,
	}
}

func (handler UserTerminalAccessRestHandlerImpl) StartTerminalSession(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var request models.UserTerminalSessionRequest
	err = decoder.Decode(&request)
	if err != nil {
		handler.Logger.Errorw("request err, StartTerminalSession", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	request.UserId = userId
	err = handler.validator.Struct(request)
	if err != nil {
		handler.Logger.Errorw("validation err, StartTerminalSession", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionCreate, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	sessionResponse, err := handler.UserTerminalAccessService.StartTerminalSession(&request)
	if err != nil {
		handler.Logger.Errorw("service err, StartTerminalSession", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, sessionResponse, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) UpdateTerminalSession(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var request models.UserTerminalSessionRequest
	err = decoder.Decode(&request)
	if err != nil {
		handler.Logger.Errorw("request err, UpdateTerminalSession", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	request.UserId = userId
	err = handler.validator.Struct(request)
	if err != nil {
		handler.Logger.Errorw("validation err, UpdateTerminalSession", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionUpdate, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	sessionResponse, err := handler.UserTerminalAccessService.UpdateTerminalSession(&request)
	if err != nil {
		handler.Logger.Errorw("service err, UpdateTerminalSession", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, sessionResponse, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) UpdateTerminalShellSession(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var request models.UserTerminalShellSessionRequest
	err = decoder.Decode(&request)
	if err != nil {
		handler.Logger.Errorw("request err, UpdateTerminalShellSession", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	err = handler.validator.Struct(request)
	if err != nil {
		handler.Logger.Errorw("validation err, UpdateTerminalShellSession", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionUpdate, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	sessionResponse, err := handler.UserTerminalAccessService.UpdateTerminalShellSession(&request)
	if err != nil {
		handler.Logger.Errorw("service err, UpdateTerminalShellSession", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, sessionResponse, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) FetchTerminalStatus(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	terminalAccessId, err := strconv.Atoi(vars["terminalAccessId"])
	if err != nil {
		handler.Logger.Errorw("request err, FetchTerminalStatus", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionGet, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	sessionResponse, err := handler.UserTerminalAccessService.FetchTerminalStatus(terminalAccessId)
	if err != nil {
		handler.Logger.Errorw("service err, FetchTerminalStatus", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, sessionResponse, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) FetchTerminalPodEvents(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	terminalAccessId, err := strconv.Atoi(vars["terminalAccessId"])
	if err != nil {
		handler.Logger.Errorw("request err, FetchTerminalPodEvents", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionGet, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}

	podEvents, err := handler.UserTerminalAccessService.FetchPodEvents(terminalAccessId)
	if err != nil {
		handler.Logger.Errorw("service err, FetchTerminalPodEvents", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, podEvents, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) FetchTerminalPodManifest(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	terminalAccessId, err := strconv.Atoi(vars["terminalAccessId"])
	if err != nil {
		handler.Logger.Errorw("request err, FetchTerminalPodManifest", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionGet, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}

	podManifest, err := handler.UserTerminalAccessService.FetchPodManifest(terminalAccessId)
	if err != nil {
		handler.Logger.Errorw("service err, FetchTerminalPodManifest", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, podManifest, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) DisconnectTerminalSession(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	terminalAccessId, err := strconv.Atoi(vars["terminalAccessId"])
	if err != nil {
		handler.Logger.Errorw("request err, DisconnectTerminalSession", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionGet, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	err = handler.UserTerminalAccessService.DisconnectTerminalSession(terminalAccessId)
	if err != nil {
		handler.Logger.Errorw("service err, DisconnectTerminalSession", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, nil, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) StopTerminalSession(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	terminalAccessId, err := strconv.Atoi(vars["terminalAccessId"])
	if err != nil {
		handler.Logger.Errorw("request err, StopTerminalSession", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionGet, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	handler.UserTerminalAccessService.StopTerminalSession(terminalAccessId)
	common.WriteJsonResp(w, nil, nil, http.StatusOK)
}

func (handler UserTerminalAccessRestHandlerImpl) DisconnectAllTerminalSessionAndRetry(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.UserService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var request models.UserTerminalSessionRequest
	err = decoder.Decode(&request)
	if err != nil {
		handler.Logger.Errorw("request err, DisconnectAllTerminalSessionAndRetry", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	request.UserId = userId
	err = handler.validator.Struct(request)
	if err != nil {
		handler.Logger.Errorw("validation err, DisconnectAllTerminalSessionAndRetry", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	if ok := handler.Enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionUpdate, "*"); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	handler.UserTerminalAccessService.DisconnectAllSessionsForUser(userId)
	sessionResponse, err := handler.UserTerminalAccessService.StartTerminalSession(&request)
	if err != nil {
		handler.Logger.Errorw("service err, DisconnectAllTerminalSessionAndRetry", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, sessionResponse, http.StatusOK)
}
