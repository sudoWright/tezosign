package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"msig/api/response"
	"msig/common/apperrors"
	"msig/models"
	"msig/repos"
	"msig/services"
	"msig/types"
	"net/http"
)

func (api *API) AuthRequest(w http.ResponseWriter, r *http.Request) {
	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	auth, err := api.provider.GetAuthProvider(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	db, err := api.provider.GetDb(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	var req models.AuthTokenReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = req.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest, err.Error()))
		return
	}

	service := services.New(repos.New(db), nil, auth, net)

	resp, err := service.AuthRequest(req)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) UnderAuthRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, isUser := ctx.Value(ContextUserKey).(types.Address)
	if !isUser || (user.Validate() != nil) {
		response.JsonError(w, apperrors.New(apperrors.ErrService))
		return
	}

	log.Print(r.Context())
	response.Json(w, "success")
}

func (api *API) Auth(w http.ResponseWriter, r *http.Request) {
	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	auth, err := api.provider.GetAuthProvider(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	db, err := api.provider.GetDb(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	var req models.SignatureReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = req.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest, err.Error()))
		return
	}

	service := services.New(repos.New(db), nil, auth, net)

	resp, err := service.Auth(req)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	//cookie := &http.Cookie{
	//	Name:     "session",
	//	Value:    "ar.EncodedCookie",
	//	Path:     "/",
	//	MaxAge:   conf.TtlCookie,
	//	Secure:   true, //api.config.IsProtocolHttps,
	//	HttpOnly: true,
	//}
	//
	//http.SetCookie(w, cookie)

	response.Json(w, resp)
}

func (api *API) RefreshAuth(w http.ResponseWriter, r *http.Request) {
	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	auth, err := api.provider.GetAuthProvider(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	db, err := api.provider.GetDb(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	var data struct {
		RefreshToken string `json:"refresh_token"`
	}

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	if data.RefreshToken == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "refresh_token"))
		return
	}

	service := services.New(repos.New(db), nil, auth, net)

	resp, err := service.RefreshAuthSession(data.RefreshToken)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	//cookie := &http.Cookie{
	//	Name:     "session",
	//	Value:    ar.EncodedCookie,
	//	Path:     "/",
	//	MaxAge:   conf.TtlCookie,
	//	Secure:   this.config.IsProtocolHttps,
	//	HttpOnly: true,
	//}
	//
	//http.SetCookie(w, cookie)

	response.Json(w, resp)
}
