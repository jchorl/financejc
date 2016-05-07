package handlers

import (
	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func GetAccountGroups(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	log.Debugf(c, "In get account groups handler, authenticated user: %+v", request.Attribute("userId"))
}
