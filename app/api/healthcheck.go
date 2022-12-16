package main

import (
	"net/http"

	"github.com/muchlist/moneymagnet/pkg/web"
)

// @Summary      Health Check
// @Description  Health Check
// @Tags         HealthCheck
// @Accept       json
// @Produce      json
// @Success      200  {object}  misc.ResponseMessage
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /healthcheck [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := web.Envelope{
		"data": "available",
	}

	err := web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
	}
}
