package cooking

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func init() {
	PublicHandlers["/login"] = handleLogin
}

func handleLogin(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		session, err := server.GetSession(writer, request)
		if err != nil {
			return
		}

		if request.Method == "POST" {
			err := request.ParseForm()
			if err != nil {
				SendError(writer, request, err, http.StatusBadRequest)
				return
			}

			username := request.PostForm.Get("username")
			password := request.PostForm.Get("password")

			user := server.GetUser(username)
			if user != EmptyUser {
				if user.Authenticate(username, password) {
					session.Values["username"] = username
					err = session.Save(request, writer)
					if err != nil {
						SendError(writer, request, err, http.StatusInternalServerError)
						return
					}
					http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
					return
				}
			}
			if server.SetMessage(writer, request, "Falscher Nutzer oder falsches Password", StatusClassCritical) != nil {
				return
			}
			logrus.Info("message set")
		}

		variables := BaseTemplateData{
			NoHeader: true,
		}
		SendTemplate(writer, request, server, "login.go.tmpl", &variables)
	}
}
