package handlers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/emicklei/go-restful"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/jchorl/financejc/transfer"
)

const autoImportPath = "../import"

func Transfer(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)
	c := appengine.NewContext(request.Request)

	files, err := ioutil.ReadDir(autoImportPath)
	if err != nil {
		log.Errorf(c, "error listing files: %+v", err)
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	for _, f := range files {
		file, err := os.Open(path.Join(autoImportPath, f.Name()))
		if err != nil {
			log.Errorf(c, "error opening file: %+v", err)
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		err = transfer.TransferQIF(c, userId, file)
		if err != nil {
			log.Errorf(c, "error importing file: %+v", err)
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
	}
}
