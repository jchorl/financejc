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

const importPath = "import"

func (s server) Transfer(request *restful.Request, response *restful.Response) {
	userId := request.Attribute("userId").(int)

	files, err := ioutil.ReadDir(importPath)
	if err != nil {
		logrus.WithField("Error", err).Error("error listing files")
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	for _, f := range files {
		logrus.WithField("File", f.Name()).Debug("about to import file")

		// skip gitkeep
		if f.Name() == ".gitkeep" {
			continue
		}

		file, err := os.Open(path.Join(importPath, f.Name()))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error": err,
				"File":  f.Name(),
			}).Error("error listing files")
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()
		logrus.WithField("File", f.Name()).Debug("file opened successfully")

		err = transfer.TransferQIF(s.Context(), userId, file)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
	}
}
