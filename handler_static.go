package cooking

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func init() {
	PublicHandlers["/static/{filename}"] = handleStatic
}

func handleStatic(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		filename, ok := mux.Vars(request)["filename"]
		if !ok {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		filepath := path.Join(server.Config.StaticDir, path.Base(filename))

		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			logrus.Errorf("cannnot open file %s: %s", filename, err)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		ext := path.Ext(filename)
		var contentType string
		switch ext {
		case ".css":
			contentType = "text/css"
		case ".js":
			contentType = "application/javascript"
		case ".jpg":
			contentType = "image/jpeg"
		case ".gif":
			contentType = "image/gif"
		case ".svg":
			contentType = "image/svg"
		case ".png":
			contentType = "image/png"
		default:
			logrus.Errorf("%s is not an allowed type", ext)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.Header().Add("Content-Type", contentType)
		writer.Write(content)
	}
}
