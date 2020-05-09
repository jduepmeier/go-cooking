package cooking

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func init() {
	PublicHandlers["/logout"] = handleLogout
}

func handleLogout(server *Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		session, err := server.SessionStore.Get(request, "session-name")
		if err != nil {
			logrus.Errorf("cannot get session: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		session.Values["username"] = ""
		err = session.Save(request, writer)
		if err != nil {
			logrus.Errorf("cannot save session: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		server.SetMessage(writer, request, "Logout erfolgreich", StatusClassSuccess)
		http.Redirect(writer, request, "/login", http.StatusTemporaryRedirect)
	}
}
