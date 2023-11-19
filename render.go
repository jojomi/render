package render

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/jojomi/asset"
)

// Renderer interface
type Renderer interface {
	ServePage(w http.ResponseWriter, req *http.Request, name string, data interface{})
	ServePageWithFuncs(w http.ResponseWriter, req *http.Request, name string, data interface{}, funcMap template.FuncMap)
	ServeError(w http.ResponseWriter, req *http.Request, code int, err error, pageData interface{})
	ErrorLogCallback(w http.ResponseWriter, req *http.Request, code int, err error) error

	GetTemplateData(name string) ([]byte, error)
	GetLayoutData(name string) ([]byte, error)

	SetDelims(left, right string)
	SetErrorLogCallback(ecb ErrorCallback)
}

// ErrorCallback type
type ErrorCallback func(w http.ResponseWriter, req *http.Request, code int, err error) error

// AppRenderer struct
type AppRenderer struct {
	TemplateDir      string
	ErrorTemplateDir string
	AssetHandler     asset.Handler
	DelimLeft        string
	DelimRight       string

	errorTemplateMap map[int]string
	errorLogCallback ErrorCallback
}

// NewAppRenderer function
func NewAppRenderer() *AppRenderer {
	return &AppRenderer{
		errorTemplateMap: make(map[int]string, 0),
	}
}

func (r *AppRenderer) SetDelims(left, right string) {
	r.DelimLeft = left
	r.DelimRight = right
}

// ServePage function
func (r *AppRenderer) ServePage(w http.ResponseWriter, req *http.Request, name string, data interface{}) {
	funcMap := template.FuncMap{
		"safeHTML": SafeHTML,
	}
	r.ServePageWithFuncs(w, req, name, data, funcMap)
}

// ServePageWithFuncs function
func (r *AppRenderer) ServePageWithFuncs(w http.ResponseWriter, req *http.Request, name string, data interface{}, funcMap template.FuncMap) {
	funcMap["safeHTML"] = SafeHTML

	tmpl, err := r.Template(name, funcMap)
	if err != nil {
		r.ServeError(w, req, http.StatusInternalServerError, err, nil)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		r.ServeError(w, req, http.StatusInternalServerError, err, nil)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func (r *AppRenderer) GetTemplateData(name string) ([]byte, error) {
	return r.AssetHandler.Get(name)
}

func (r *AppRenderer) GetLayoutData(name string) ([]byte, error) {
	return r.AssetHandler.Get(name)
}

// Template function
func (r *AppRenderer) Template(name string, funcMap template.FuncMap) (*template.Template, error) {
	data, err := r.GetTemplateData(name)
	if err != nil {
		return nil, err
	}
	tmpl := template.New(name)
	if r.DelimLeft != "" || r.DelimRight != "" {
		tmpl = tmpl.Delims(r.DelimLeft, r.DelimRight)
	}
	return tmpl.Funcs(funcMap).Parse(string(data))
}

// SetErrorLogCallback function
func (r *AppRenderer) SetErrorLogCallback(ecb ErrorCallback) {
	r.errorLogCallback = ecb
}

// ErrorLogCallback function
func (r *AppRenderer) ErrorLogCallback(w http.ResponseWriter, req *http.Request, code int, err error) error {
	if r.errorLogCallback == nil {
		return nil
	}
	return r.errorLogCallback(w, req, code, err)
}

var SafeHTML = func(str string) template.HTML {
	return template.HTML(str)
}
