package apperrors

import (
	"fmt"
	"net/http"
)

type (
	ErrCode string
)

// common
const (
	ErrService             ErrCode = "ERR_SERVICE"
	ErrNotFound            ErrCode = "ERR_NOT_FOUND"
	ErrAlreadyExists       ErrCode = "ERR_ALREADY_EXISTS"
	ErrBadRequest          ErrCode = "ERR_BAD_REQUEST"
	ErrBadParam            ErrCode = "ERR_BAD_PARAM"
	ErrNotAllowed          ErrCode = "ERR_NOT_ALLOWED"
	ErrBadJwt              ErrCode = "ERR_BAD_JWT"
	ErrBadAuth             ErrCode = "ERR_BAD_AUTH"
	ErrBadSignature        ErrCode = "ERR_BAD_SIGNATURE"
	ErrBadAuthCookie       ErrCode = "ERR_BAD_AUTH_COOKIE"
	ErrNotEnoughPermission ErrCode = "ERR_NOT_ENOUGH_PERMISSION"

	ErrUserAlreadyVerified ErrCode = "ERR_ALREADY_VERYFIED"
)

type (
	ServiceError interface {
		error
		ErrorCode() ErrCode
		ToMap(http.ResponseWriter) map[string]interface{}
		GetHttpCode() int
	}

	Error struct {
		Code        ErrCode `json:"code"`
		Value       string  `json:"value,omitempty"`
		Description string  `json:"description,omitempty"`
	}
)

func (e Error) Error() string {
	return fmt.Sprintf("%s %s", string(e.Code), e.Value)
}

func (e Error) ErrorCode() ErrCode {
	return e.Code
}

// ToMap converts Error object to map[string]interface{}
func (e Error) ToMap() map[string]interface{} {
	r := map[string]interface{}{
		"error": string(e.Code),
	}

	if e.Value != "" {
		r["value"] = e.Value
	}

	if e.Description != "" {
		r["description"] = e.Description
	}

	return r
}

// GetHttpCode return a Http error code
func (e Error) GetHttpCode() int {
	switch e.Code {
	case ErrService:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// New creates an Error object
func New(code ErrCode, value ...string) *Error {
	e := &Error{Code: code}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// NewWithDesc creates an Error object with description
func NewWithDesc(code ErrCode, desc string, value ...string) *Error {
	e := &Error{Code: code, Description: desc}
	if len(value) > 0 {
		e.Value = value[0]
	}
	return e
}

// FromError creates a new Error (ErrService) from common golang error
func FromError(err error) *Error {
	if err != nil {
		return &Error{
			Code:  ErrService,
			Value: "",
		}
	}

	return nil
}
