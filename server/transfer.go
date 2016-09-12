package server

import (
	"github.com/jchorl/financejc/server/transfer"

	"github.com/emicklei/go-restful"
)

func (s server) Transfer(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)
	err := transfer.AutoImport(s.ContextWithUser(userId))
	if err != nil {
		writeError(response, err)
		return
	}
}
