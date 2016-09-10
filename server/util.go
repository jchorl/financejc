package server

import (
	"fmt"

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
	response.WriteEntity(entity.Values())
}
