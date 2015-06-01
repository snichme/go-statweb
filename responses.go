package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	NotFound    = Template(404, "404.html", nil)
	ServerError = Empty(500)
)
var funcMap = template.FuncMap{
	"label": strings.Title,
	"menu": func(s string) map[string]string {
		links := make(map[string]string)
		files, err := ioutil.ReadDir("page/" + s)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			var filename = f.Name()
			var extension = filepath.Ext(filename)
			var name = filename[0 : len(filename)-len(extension)]
			links[name] = strings.Title(name)
		}
		return links
	},
}

type Response interface {
	WriteTo(out http.ResponseWriter)
}

type NormalResponse struct {
	status int
	body   []byte
	header http.Header
}

func (r *NormalResponse) WriteTo(out http.ResponseWriter) {
	header := out.Header()
	for k, v := range r.header {
		header[k] = v
	}
	out.WriteHeader(r.status)
	out.Write(r.body)
}

func (r *NormalResponse) Cache(ttl string) *NormalResponse {
	return r.Header("Cache-Control", "public,max-age="+ttl)
}

func (r *NormalResponse) Header(key, value string) *NormalResponse {
	r.header.Set(key, value)
	return r
}

type TemplateResponse struct {
	status int
	data   interface{}
	header http.Header

	template string
}

func (r *TemplateResponse) WriteTo(out http.ResponseWriter) {
	templatePath := fmt.Sprintf("layout/%s", r.template)
	t := template.Must(template.New(r.template).
		Funcs(funcMap).
		ParseFiles(templatePath))
	t.Execute(out, r.data)
}

func (r *TemplateResponse) Cache(ttl string) *TemplateResponse {
	return r.Header("Cache-Control", "public,max-age="+ttl)
}

func (r *TemplateResponse) Header(key, value string) *TemplateResponse {
	r.header.Set(key, value)
	return r
}

func Empty(status int) *NormalResponse {
	return Respond(status, nil)
}

func Json(status int, body interface{}) *NormalResponse {
	var b []byte
	var err error
	if b, err = json.Marshal(body); err != nil {
		return Error("body json marshal", err)
	}
	return Respond(status, b).Header("Content-Type", "application/json")
}

func Text(status int, body string) *NormalResponse {
	return Respond(status, []byte(body)).Header("Content-Type", "text/plain")
}
func Template(status int, template string, data interface{}) *TemplateResponse {
	return &TemplateResponse{
		status:   status,
		header:   make(http.Header),
		template: template,
		data:     data,
	}
}
func Error(message string, err error) *NormalResponse {
	log.Println(message, err)
	return ServerError
}

func Respond(status int, body []byte) *NormalResponse {
	return &NormalResponse{
		body:   body,
		status: status,
		header: make(http.Header),
	}
}
