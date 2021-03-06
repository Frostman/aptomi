package api

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/lang/yaml"
	"github.com/Aptomi/aptomi/pkg/object"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
)

type handler = func(*http.Request, httprouter.Params) Response
type streamHandler = func(http.ResponseWriter, *http.Request, httprouter.Params)

func (a *api) get(path string, handler handler) {
	a.handle("GET", path, handler)
}

func (a *api) post(path string, handler handler) {
	a.handle("POST", path, handler)
}

func (a *api) handle(method string, path string, handler handler) {
	a.router.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		resp := handler(request, params)
		if resp != nil {
			write(writer, resp)
		} else {
			writer.WriteHeader(http.StatusNoContent)
		}
	})
}

func (a *api) getStream(path string, handler streamHandler) {
	a.router.GET(path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		handler(writer, request, params)
	})
}

func (a *api) read(request *http.Request) []object.Base {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Panicf("Error while reading bytes from request Body: %s", err)
	}

	objects, err := a.codec.UnmarshalOneOrMany(body)
	if err != nil {
		// todo response with some bad request status code
		log.Panicf("Error unmarshalling policy update request: %s", err)
	}

	return objects
}

func write(writer io.Writer, resp Response) {
	data := yaml.SerializeObject(resp)
	_, wErr := fmt.Fprint(writer, data)
	if wErr != nil {
		log.Panicf("Error while writing response: %s", wErr)
	}
}
