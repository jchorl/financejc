package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/jchorl/financejc/server/transfer"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

const autoImportPath = "../import"

func (s server) Transfer(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(string)

	files, err := ioutil.ReadDir(autoImportPath)
	if err != nil {
		logrus.WithField("Error", err).Error("error listing files")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	for _, f := range files {
		file, err := os.Open(path.Join(autoImportPath, f.Name()))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error": err,
				"File":  f.Name(),
			}).Error("error listing files")
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		err = transfer.TransferQIF(s.Context(), userId, file)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
	}
}
