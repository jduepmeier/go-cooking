package cooking

import (
	"fmt"
	"net/http"
)

func init() {
	Handlers["/add"] = handleGet
}

// AddTemplateData contains the template data for the add template.
type AddTemplateData struct {
	BaseTemplateData
	Link   string
	Recipe Recipe
}

func handleGet(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "POST":
			err := request.ParseForm()
			if err != nil {
				SendError(writer, request, fmt.Errorf("cannot parse from: %s", err), http.StatusBadRequest)
				return
			}
			recipe := Recipe{
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
			err = server.Storage.Add(recipe)
			if err != nil {
				SendError(writer, request, err, http.StatusBadRequest)
				return
			}
			writer.Header().Add("Location", "/?success=\"added successfully\"")
			writer.WriteHeader(http.StatusTemporaryRedirect)
		case "GET":
			SendTemplate(writer, request, server, "add.go.tmpl", &AddTemplateData{
				Link: "/add",
			})
		}
	}
}
