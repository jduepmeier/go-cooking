package cooking

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	Handlers["/delete/{id:[0-9]+}"] = handleDelete
}

func handleDelete(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		arguments := mux.Vars(request)
		id, err := strconv.ParseUint(arguments["id"], 10, 64)
		if err != nil {
			SendError(writer, request, err, http.StatusBadRequest)
		}
		switch request.Method {
		case "POST":
			err = server.Storage.Delete(id)
			if err != nil {
				SendError(writer, request, err, http.StatusBadRequest)
				return
			}
			server.SetMessage(writer, request, "deleted successfully", StatusClassSuccess)
			http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		case "GET":
			recipe, err := server.Storage.Get(id)
			if err != nil {
				SendError(writer, request, err, http.StatusNotFound)
				return
			}

			SendTemplate(writer, request, server, "delete.go.tmpl", &AddTemplateData{
				Recipe: recipe,
			})
		}
	}
}
