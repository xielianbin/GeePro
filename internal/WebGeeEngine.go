package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type WebGeeEngine struct {
	*RouterGroup
	groups []*RouterGroup
	router IRouter //定义一个路由映射表，将对于路径的请求和处理函数存储起来

	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

// 定义构造函数
func NewWebGeeEngine(webGeeRouter *WebGeeRouter) *WebGeeEngine {
	// 当声明一个 map、slice 或 channel 类型的变量时，它只是声明了一个变量，并没有分配内存来存储数据。因此，在使用之前，必须初始化它们。
	// 在 NewWebGee 构造函数中，使用 make(map[string]HandlerFunc) 来初始化 router 字段。
	// 这确保了当 WebGee 实例被创建时，router 字段已经是一个可以存储数据的 map。
	router := webGeeRouter
	engine := &WebGeeEngine{router: router}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// GetHtmlTemplates returns the parsed HTML templates.
func (wg *WebGeeEngine) GetHtmlTemplates() *template.Template {
	return wg.htmlTemplates
}

func (wg *WebGeeEngine) SetFuncMap(funcMap template.FuncMap) {
	wg.funcMap = funcMap
}

func (wg *WebGeeEngine) LoadHTMLGlob(pattern string) {
	wg.htmlTemplates = template.Must(template.New("").Funcs(wg.funcMap).ParseGlob(pattern))
}
func (wg *WebGeeEngine) addRoute(method string, pattern string, handler HandlerFunc) {

	wg.router.AddRouter(method, pattern, handler)

}
func (wg *WebGeeEngine) GET(pattern string, handler HandlerFunc) {
	wg.addRoute("GET", pattern, handler)
}
func (wg *WebGeeEngine) POST(pattern string, handler HandlerFunc) {
	wg.addRoute("POST", pattern, handler)
}
func (wg *WebGeeEngine) RUN(addr string) (err error) {
	fmt.Println("开启监听")
	return http.ListenAndServe(addr, wg)
}

// 在开启监听的时候，需要放入地址和一个接口，这个接口必须实现serverHttp方法，当检测到请求的时候就会调用这个方法
func (wg *WebGeeEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range wg.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewWebGeeContext(w, req)
	c.SetHandlers(middlewares)
	c.SetEngine(wg)
	fmt.Println("设置引擎")
	//fmt.Println(wg.router)
	if wg == nil {
		fmt.Println("引擎为空")
		return
	}
	if wg.router == nil {
		fmt.Println("路由为空")
		return
	}
	wg.router.Handle(c)
}
