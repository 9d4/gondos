package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"gondos/internal/app"
)

func (si serverImpl) deliverErr(w http.ResponseWriter, r *http.Request, err error) {
	log.Debug().Caller().Err(err).Send()
	var (
		appUserErr      app.UserError
		appValidateErr  app.ValidationError
		appValidateErrs app.ValidationErrors
	)

	response := Error{
		Code:    "internal",
		Message: "Internal server error",
	}

	deliverCustom := func(status int, data any) {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(data)
	}
	deliver := func(status int) {
		deliverCustom(status, response)
	}

	switch {
	case errors.As(err, &appUserErr):
		response.Code = appUserErr.Code()
		response.Message = appUserErr.Message()
		deliver(getStatusFromKind(appUserErr.Kind()))
		return

	case errors.Is(err, os.ErrPermission):
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
		return

	case errors.As(err, &appValidateErr):
		response := ValidationError{
			Code:    "validation",
			Message: "Your request didn't validate. Please check your input and try again.",
		}
		response.Params = parseValidationErrorParams(appValidateErr)

		deliverCustom(http.StatusUnprocessableEntity, response)
		return

	case errors.As(err, &appValidateErrs):
		response := ValidationError{
			Code:    "validation",
			Message: "Your request didn't validate. Please check your input and try again.",
		}

		for _, ve := range appValidateErrs {
			response.Params = append(response.Params, parseValidationErrorParams(ve)...)
		}

		deliverCustom(http.StatusUnprocessableEntity, response)
		return
	}

	deliver(http.StatusInternalServerError)
}

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

var appErrKindStatus = map[app.ErrorKind]int{
	app.DuplicateErrorKind:  http.StatusConflict,
	app.ValidationErrorKind: http.StatusUnprocessableEntity,
}

func getStatusFromKind(kind app.ErrorKind) int {
	if i, ok := appErrKindStatus[kind]; ok {
		return i
	}
	return http.StatusInternalServerError
}
