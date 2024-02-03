package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog/log"

	"gondos/internal/app"
)

// deliverErr parses errors and responds to client
func (si serverImpl) deliverErr(w http.ResponseWriter, r *http.Request, err error) {
	log.Debug().Caller().Err(err).Send()
	var (
		errJsonSyntax    *json.SyntaxError
		errJsonMarshaler *json.MarshalerError
		errAppUser       *app.UserError
		errAppValidate   app.ValidationError
		errAppValidates  app.ValidationErrors
	)

	res := Error{
		Code:    "internal",
		Message: "Internal server error",
	}

	sendRes := func(status int) {
		sendJSON(w, status, res)
	}

	switch {
	case errors.Is(err, io.EOF):
		res.Code = "input.required"
		res.Message = "No request body"
		sendRes(http.StatusBadRequest)
		return

	case isJWTAuthError(err):
		res.Code = "unauthorized"
		res.Message = "Authentication required"
		sendRes(http.StatusUnauthorized)
		return

	case errors.As(err, &errJsonSyntax):
		res.Code = "json.syntax"
		res.Message = "Json syntax error"
		sendRes(http.StatusBadRequest)
		return

	case errors.As(err, &errJsonMarshaler):
		res.Code = "json.parse"
		res.Message = "Cannot parse json"
		sendRes(http.StatusBadRequest)
		return

	case errors.As(err, &errAppUser):
		res.Code = errAppUser.Code()
		res.Message = errAppUser.Message()
		sendRes(getStatusFromKind(errAppUser.Kind()))
		return

	case errors.Is(err, os.ErrPermission):
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
		return

	case errors.As(err, &errAppValidate):
		response := ValidationError{
			Code:    "validation",
			Message: "Your request didn't validate. Please check your input and try again.",
		}
		response.Params = parseValidationErrorParams(errAppValidate)

		sendJSON(w, http.StatusUnprocessableEntity, response)
		return

	case errors.As(err, &errAppValidates):
		response := ValidationError{
			Code:    "validation",
			Message: "Your request didn't validate. Please check your input and try again.",
		}

		for _, ve := range errAppValidates {
			response.Params = append(response.Params, parseValidationErrorParams(ve)...)
		}

		sendJSON(w, http.StatusUnprocessableEntity, response)
		return
	}

	sendRes(http.StatusInternalServerError)
}

var appErrKindStatus = map[app.ErrorKind]int{
	app.ErrorKindBad:        http.StatusBadRequest,
	app.ErrorKindNotFound:   http.StatusNotFound,
	app.ErrorKindDuplicate:  http.StatusConflict,
	app.ErrorKindValidation: http.StatusUnprocessableEntity,
}

func getStatusFromKind(kind app.ErrorKind) int {
	if i, ok := appErrKindStatus[kind]; ok {
		return i
	}
	return http.StatusInternalServerError
}

// parseValidationErrorParams converts app.ValidationError to slice of ValidationErrorParams
func parseValidationErrorParams(err app.ValidationError) []ValidationErrorParams {
	params := []ValidationErrorParams{}

	for _, ve := range err.ValidationErrors {
		p := map[string]interface{}{
			"field": err.FieldName(),
			"tag":   ve.Tag(),
			"param": ve.Param(),
		}
		params = append(params, p)
	}
	return params
}

func isJWTAuthError(err error) bool {
	return errors.Is(err, jwtauth.ErrUnauthorized) ||
		errors.Is(err, jwtauth.ErrExpired) ||
		errors.Is(err, jwtauth.ErrNBFInvalid) ||
		errors.Is(err, jwtauth.ErrIATInvalid) ||
		errors.Is(err, jwtauth.ErrNoTokenFound) ||
		errors.Is(err, jwtauth.ErrAlgoInvalid)
}
