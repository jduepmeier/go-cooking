// Package cooking contains all the cooking server code.
package cooking

import (
	"fmt"
	"html"
	"net/http"
	"path"
	"strings"
	"text/template"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Handlers contains all known http handlers.
	// This will be filled with init functions.
	Handlers = make(map[string]HandleFunc)
	// PublicHandlers are handler which will not have an authentication check.
	// this will be filled with init functions.
	PublicHandlers = make(map[string]HandleFunc)
)

var (
	// EmptyUser is the empty user
	EmptyUser = User{}
)

// HandleFunc defines a function which returns a http.HandleFunc
type HandleFunc = func(server *Server) http.HandlerFunc

// Config contains the configuration for the server.
type Config struct {
	// Addr is defined as <bind>:<port>
	Addr string `yaml:"addr"`
	// Templates is the path to the templates folder.
	Templates   string `yaml:"templates"`
	StoragePath string `yaml:"storage-path"`
	StaticDir   string `yaml:"static-dir"`
}

// User contains
type User struct {
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	authenticated bool   `yaml:"-"`
}

const (
	// StatusClassSuccess should be returned if the message is positive.
	StatusClassSuccess = "toast success"
	// StatusClassWarning should be returned if the message is not critical.
	StatusClassWarning = "toast warning"
	// StatusClassCritical should be returned if the message is critical.
	StatusClassCritical = "toast error"
)

// StatusMessage is a message to return to the requester via template.
type StatusMessage struct {
	Message string
	Class   string
}

// HTML returns the html encoded status message.
func (message StatusMessage) HTML() string {
	msg := ""
	if message.Message != "" {
		msg = fmt.Sprintf(`<div class=%q>%s</div>`, message.Class, html.EscapeString(message.Message))
	}
	return msg
}

// Authenticate authenticate a user in
func (user *User) Authenticate(username, password string) bool {
	if user.Username == username &&
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
		return true
	}
	return false
}

// GetUser returns the user with
func (server *Server) GetUser(username string) User {
	user, err := server.Storage.GetUser(username)
	if err != nil {
		logrus.Errorf("cannot get user from database: %s", err)
		return EmptyUser
	}
	return user
}

// Server is the server.
type Server struct {
	Config       *Config
	Templates    *template.Template
	Storage      *Storage
	SessionStore *sessions.CookieStore
	done         chan struct{}
}

// DefaultConfig returns a config struct with default values.
func DefaultConfig() *Config {
	return &Config{
		Addr:        "localhost:8080",
		Templates:   "templates",
		StoragePath: "test.db3",
		StaticDir:   "static",
	}
}

// NewServer creates a new server with the given configuration.
func NewServer(config *Config) (*Server, error) {
	server := &Server{
		Config: config,
	}
	server.Storage = &Storage{
		Path: server.Config.StoragePath,
	}
	err := server.Storage.Connect()
	if err != nil {
		return server, fmt.Errorf("cannot connect to storage: %s", err)
	}

	return server, nil
}

func (server *Server) loadTemplates() (err error) {
	logrus.Infof("load templates")
	globPattern := path.Join(server.Config.Templates, "*.tmpl")
	templates, err := template.ParseGlob(globPattern)
	if err == nil {
		server.Templates = templates
	}
	return err
}

func (server *Server) watchTemplates() (err error) {
	server.done = make(chan struct{})
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logrus.Debugf("got event %s", event.String())
				if strings.HasSuffix(event.Name, ".go.tmpl") {
					err = server.loadTemplates()
					if err != nil {
						logrus.Errorf("cannot reload templates: %s", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Errorf("got error on watching templates: %s", err)
			}
		}
	}()
	logrus.Infof("begin watching %s", server.Config.Templates)
	err = watcher.Add(server.Config.Templates)
	if err != nil {
		return err
	}
	<-server.done
	return nil
}

// Serve creates a http socket and serves the server.
func (server *Server) Serve() error {
	server.SessionStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
	go server.watchTemplates()
	err := server.loadTemplates()
	if err != nil {
		return err
	}
	mux := mux.NewRouter()
	for pattern, handleFunc := range Handlers {
		mux.Handle(pattern, server.sessionHandler(handleFunc(server)))
	}
	for pattern, handleFunc := range PublicHandlers {
		mux.Handle(pattern, handleFunc(server))
	}
	err = http.ListenAndServe(server.Config.Addr, mux)
	logrus.Infof("Got error %s, shutting down...", err)
	server.done <- struct{}{}

	return nil
}

func (server *Server) sessionHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		session, err := server.GetSession(writer, request)
		if err != nil {
			return
		}
		user, ok := session.Values["username"]
		if ok {
			logrus.Info("user %s found", user)
			if userData, ok := user.(string); ok {
				if userData != "" {
					logrus.Info("user %v could be parsed", userData)
					handler(writer, request)
					return
				}
			}
		}

		http.Redirect(writer, request, "/login", http.StatusTemporaryRedirect)
	}
}

// GetSession returns the user session object.
func (server *Server) GetSession(writer http.ResponseWriter, request *http.Request) (*sessions.Session, error) {
	session, err := server.SessionStore.Get(request, "session-name")
	if err != nil {
		logrus.Errorf("cannot get session: %s", err)
		err = session.Save(request, writer)
		if err != nil {
			logrus.Errorf("cannot save session: %s", err)
		}
	}
	return session, nil
}

// GetMessage returns the message from the session if exists.
func (server *Server) GetMessage(writer http.ResponseWriter, request *http.Request) (message StatusMessage, err error) {
	session, err := server.GetSession(writer, request)
	if err != nil {
		return
	}

	messages := session.Flashes()
	if len(messages) > 0 {
		message = messages[0].(StatusMessage)
	}
	return
}

// SetMessage sets the message for the next template call.
func (server *Server) SetMessage(writer http.ResponseWriter, request *http.Request, message string, messageClass string) error {
	session, err := server.GetSession(writer, request)
	if err != nil {
		return err
	}
	session.AddFlash(StatusMessage{
		Message: message,
		Class:   messageClass,
	})

	return nil
}

// HashPassword hashes a users password
func (server *Server) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 5)
}
