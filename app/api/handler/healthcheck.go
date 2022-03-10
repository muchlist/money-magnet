package handler

import (
	"net/http"

	"github.com/muchlist/moneymagnet/foundation/web"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	env := web.Envelope{
		"status": "available",
	}

	err := web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
	}
}
