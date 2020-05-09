package cooking

import (
	"net/http"
)

func init() {
	Handlers["/"] = handleRoot
}

// RootTemplateData contains the template data for the root template.
type RootTemplateData struct {
	BaseTemplateData
	Recipes []Recipe
}

func handleRoot(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		recipes, err := server.Storage.GetRecipes()
		if err != nil {
			SendError(writer, request, err, http.StatusInternalServerError)
			return
		}

		variables := RootTemplateData{
			Recipes: recipes,
		}
		SendTemplate(writer, request, server, "root.go.tmpl", &variables)
	}
}
