package render

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jojomi/asset"
)

// Renderer interface
type Renderer interface {
	ServePage(w http.ResponseWriter, req *http.Request, name string, data interface{})
	ServeError(w http.ResponseWriter, req *http.Request, code int, err error, pageData interface{})
	ErrorLogCallback(w http.ResponseWriter, req *http.Request, code int, err error) error

	GetTemplateData(name string) ([]byte, error)
	GetLayoutData(name string) ([]byte, error)

	SetErrorLogCallback(ecb ErrorCallback)
}

// ErrorCallback type
type ErrorCallback func(w http.ResponseWriter, req *http.Request, code int, err error) error

// AppRenderer struct
type AppRenderer struct {
	TemplateDir      string
	ErrorTemplateDir string
	AssetHandler     asset.Handler

	errorTemplateMap map[int]string
	errorLogCallback ErrorCallback
}

// NewAppRenderer function
func NewAppRenderer() *AppRenderer {
	return &AppRenderer{
		errorTemplateMap: make(map[int]string, 0),
	}
}

// ServePage function
func (r *AppRenderer) ServePage(w http.ResponseWriter, req *http.Request, name string, data interface{}) {
	tmpl, err := r.Template(name)
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
func (r *AppRenderer) Template(name string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"StringsJoin": strings.Join,
		"SHAFile": func(outputFilename, fsFilename string) string {
			hasher := sha256.New()
			s, err := ioutil.ReadFile(fsFilename)
			hasher.Write(s)
			if err != nil {
				return outputFilename // error case, fallback to filename only, file might not exist
			}
			return outputFilename + "?" + hex.EncodeToString(hasher.Sum(nil))
		},
		"safeHTML": func(str string) template.HTML {
			return template.HTML(str)
		},
	}
	data, err := r.GetTemplateData(name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Funcs(funcMap).Parse(string(data))
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
