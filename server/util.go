package server

import (
	"fmt"
	"net/http"

	"github.com/jchorl/financejc/constants"

	"github.com/emicklei/go-restful"
)

type Paginated interface {
	Next() string
	Values() []interface{}
}

func writePaginatedEntity(request *restful.Request, response *restful.Response, entity Paginated) {
	u := request.Request.URL
	u.RawQuery = ""
	response.AddHeader("Link", fmt.Sprintf("<%s?start=%s>; rel=\"next\"", u, entity.Next()))
	if entity.Values() != nil {
		response.WriteEntity(entity.Values())
		return
	}

	response.WriteEntity([]interface{}{})
}

func writeError(response *restful.Response, err error) {
	switch err {
	case constants.NotLoggedIn:
		response.WriteError(http.StatusUnauthorized, err)
	case constants.Forbidden:
		response.WriteError(http.StatusForbidden, err)
	case constants.BadRequest:
		response.WriteError(http.StatusBadRequest, err)
	default:
		response.WriteError(http.StatusInternalServerError, err)
	}
}
