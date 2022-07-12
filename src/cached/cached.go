package cached

import (
	"net/http"
)

const (
	HeaderContentType = "Content-Type"
)

//var hlog = logrus.WithFields("component")

type Entry struct {
	ContentType string
	Content []byte
}

type Handler struct {
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//e := Entry{}


}

func (h *Handler) okResponse(writer http.ResponseWriter, e Entry) {
	writer.Header().Set(HeaderContentType, e.ContentType)
	//writer.Write()
}
