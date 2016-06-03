package handlers

import (
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"upload"
)

func Upload(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)
	c := appengine.NewContext(request.Request)

	file, err := os.Open("/home/josh/downloads/transactions.txt")
	if err != nil {
		log.Errorf(c, "error opening file: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	err = upload.Upload(c, userId, file, "TSV")
	if err != nil {
		log.Errorf(c, "error importing file: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
}
