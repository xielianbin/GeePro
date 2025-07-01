package internal

import (
	"html/template"
	"net/http"
)

type IEngin interface {
	RUN(addr string) error
	GET(pattern string, handler HandlerFunc)
	POST(pattern string, handler HandlerFunc)
	ServeHTTP(w http.ResponseWriter, req *http.Request)

	//SetRouter(router *IRouter)
	//GetRouter() *IRouter
	SetFuncMap(funcMap template.FuncMap)
	LoadHTMLGlob(pattern string)
	GetHtmlTemplates() *template.Template
}
