package cooking

import (
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	validFilenameCharacters = regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`)
)

func init() {
	PublicHandlers["/static/{filename}"] = handleStatic
	PublicHandlers["/favicon.ico"] = handleStatic
}

func sanitizeFilename(filename string) (string, bool) {
	filename = strings.Replace(filename, "\n", "", -1)
	filename = strings.Replace(filename, "\r", "", -1)

	return filename, validFilenameCharacters.MatchString(filename)
}

func handleStatic(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var filename string
		var ok bool
		if strings.HasSuffix(request.URL.Path, "favicon.ico") {
			filename = "favicon.ico"
		} else {
			filename, ok = mux.Vars(request)["filename"]
			if !ok {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
		}

		filename, ok = sanitizeFilename(filename)
		if !ok {
			logrus.Debugf("filename %q does not match allowed pattern", filename)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		filename = path.Base(filename)
		if filename == "." || filename == "/" {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		filepath := path.Join(server.Config.StaticDir, filename)

		content, err := os.ReadFile(filepath)
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
			contentType = "image/svg+xml"
		case ".png":
			contentType = "image/png"
		case ".ico":
			contentType = "image/x-icon"
		default:
			logrus.Errorf("%s is not an allowed type", ext)
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		writer.Header().Add("Content-Type", contentType)
		writer.Write(content)
	}
}
