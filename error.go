package render

import (
	"fmt"
	"html/template"
	"net/http"

	raven "github.com/getsentry/raven-go"
)

// ServeError function
func (a AppRenderer) ServeError(w http.ResponseWriter, req *http.Request, code int, err error, pageData interface{}) {
	fmt.Println("serving error", err, code, pageData)
	message := http.StatusText(code)

	a.ErrorLogCallback(w, req, code, err)

	http.Error(w, message, code)

	// TODO panic wenn 500er? oder alle?
	// panic(err)
}

// SentryLogger var for convenience
var SentryLogger = func(w http.ResponseWriter, req *http.Request, code int, err error) error {
	// 401, 403, and 404 are not relevant for sentry logging
	if code == http.StatusUnauthorized || code == http.StatusForbidden || code == http.StatusNotFound {
		return nil
	}

	// build error artificially
	if err == nil {
		err = fmt.Errorf("HTTP error %d", code)
	}
	raven.SetHttpContext(raven.NewHttp(req))
	raven.CaptureError(err, nil)
	return nil
}

func getErrorTemplate(code int) (*template.Template, error) {
	return nil, nil
}
