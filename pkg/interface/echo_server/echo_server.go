package echoserver

import (
	"net/http"

	"github.com/taroxii/vote-api/pkg/entity"
)

type ResponseError struct {
	Message string `json:"message"`
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case entity.ErrInternalServerError:
		return http.StatusInternalServerError
	case entity.ErrNotFound:
		return http.StatusNotFound
	case entity.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
