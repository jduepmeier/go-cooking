package cooking

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	Handlers["/edit/{id:[0-9]+}"] = handleEdit
}

func handleEdit(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		arguments := mux.Vars(request)
		id, err := strconv.ParseUint(arguments["id"], 10, 64)
		if err != nil {
			SendError(writer, request, err, http.StatusBadRequest)
		}
		switch request.Method {
		case "POST":
			err := request.ParseForm()
			if err != nil {
				SendError(writer, request, fmt.Errorf("cannot parse from: %s", err), http.StatusBadRequest)
				return
			}
			recipe := Recipe{
				ID:     id,
				Name:   request.PostForm.Get("name"),
				Length: request.PostForm.Get("length"),
				Source: request.PostForm.Get("source"),
			}
			freshness := request.PostForm.Get("freshness")
			err = recipe.ParseFreshness(freshness)
			if err != nil {
				SendError(writer, request, err, http.StatusBadRequest)
				return
			}
			err = recipe.Validate()
			if err != nil {
				SendError(writer, request, err, http.StatusBadRequest)
				return
			}
			err = server.Storage.Update(recipe)
			if err != nil {
				SendError(writer, request, err, http.StatusBadRequest)
				return
			}
			writer.Header().Add("Location", "/?success=\"updated successfully\"")
			writer.WriteHeader(http.StatusTemporaryRedirect)
		case "GET":
			recipe, err := server.Storage.Get(id)
			if err != nil {
				SendError(writer, request, err, http.StatusBadRequest)
				return
			}

			SendTemplate(writer, request, server, "add.go.tmpl", &AddTemplateData{
				Link:   fmt.Sprintf("/edit/%d", id),
				Recipe: recipe,
			})
		}
	}
}
