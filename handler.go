package cooking

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// TemplateData contains the map for templates data.
type TemplateData interface {
	Message() StatusMessage
	SetMessage(message StatusMessage)
}

// BaseTemplateData is the base implementation for template data.
type BaseTemplateData struct {
	message  StatusMessage
	NoHeader bool
	Printer  bool
}

// Message returns the message saved.
func (data *BaseTemplateData) Message() StatusMessage {
	return data.message
}

// SetMessage sets the message.
func (data *BaseTemplateData) SetMessage(message StatusMessage) {
	data.message = message
}

// SendError over http.
func SendError(writer http.ResponseWriter, request *http.Request, err error, errCode int) {
	logrus.Errorf("send error: %d %s", errCode, err)
	writer.WriteHeader(errCode)
	writer.Write([]byte(err.Error()))
}

// SendTemplate sends the given template.
func SendTemplate(writer http.ResponseWriter, request *http.Request, server *Server, templateName string, data TemplateData) {
	message, err := server.GetMessage(writer, request)
	if err != nil {
		return
	}
	data.SetMessage(message)
	logrus.Infof("message is %s", message.HTML())
	writer.Header().Add("Content-Type", "text/html")
	err = server.Templates.ExecuteTemplate(writer, templateName, data)
	if err != nil {
		SendError(writer, request, err, http.StatusInternalServerError)
	}
}
