package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{} //给json数据取别名
// 设计上下文(Context)，封装 Request 和 Response ，提供对 JSON、HTML 等返回类型的支持。
type WebGeeContext struct {
	// 目标
	writer http.ResponseWriter
	req    *http.Request
	// 请求信息
	path   string
	method string
	// 响应信息
	statusCode int
	params     map[string]string //将解析的参数存到这个里面
	// 中间件
	handlers []HandlerFunc
	index    int
	engine   IEngin //支持根据模板文件名选择模板进行渲染。
}

func (c *WebGeeContext) SetEngine(engine IEngin) {
	c.engine = engine
}
func (c *WebGeeContext) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
func (c *WebGeeContext) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.GetHtmlTemplates().ExecuteTemplate(c.GetWriter(), name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
func (c *WebGeeContext) GetHandlers() []HandlerFunc {
	return c.handlers
}

func (c *WebGeeContext) SetHandlers(h []HandlerFunc) {
	c.handlers = h
}
func (c *WebGeeContext) SetParam(params map[string]string) {
	c.params = params
}

func (c *WebGeeContext) Param(key string) string {
	fmt.Println(key)
	value, _ := c.params[key]
	return value
}
func (c *WebGeeContext) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
func NewWebGeeContext(w http.ResponseWriter, req *http.Request) IContext {
	var c IContext
	c = &WebGeeContext{
		writer: w,
		req:    req,
		path:   req.URL.Path,
		method: req.Method,
		index:  -1,
	}
	return c
}
func (c *WebGeeContext) PostForm(key string) string {
	//FormValue方法‌：这是Go语言net/http包中Request对象的方法，会自动解析URL查询参数和POST表单数据（包括application/x-www-form-urlencoded和multipart/form-data编码），返回第一个匹配键名的字符串值
	return c.req.FormValue(key)
}
func (c *WebGeeContext) Query(key string) string {
	//URL.Query()会解析URL中?后的查询字符串（如?name=John&age=20），返回url.Values类型（本质是map[string][]string）89
	//Get(key)方法从映射中提取首个匹配键的值（若键不存在返回空字符串）
	// 处理GET请求 /search?q=golang
	//searchTerm := c.Req.URL.Query().Get("q")
	return c.req.URL.Query().Get(key)
}
func (c *WebGeeContext) Status(code int) {
	c.statusCode = code
	//设置响应状态码
	c.writer.WriteHeader(code)
}
func (c *WebGeeContext) SetHeader(key string, value string) {
	//设置响应头的键值对
	c.writer.Header().Set(key, value)
}
func (c *WebGeeContext) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	//直接写入字符串响应
	c.writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *WebGeeContext) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	//将变长参数values序列化为JSON输出
	encoder := json.NewEncoder(c.writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.writer, err.Error(), 500)
	}
}
func (c *WebGeeContext) Data(code int, data []byte) {
	c.Status(code)
	c.writer.Write(data)
}

func (c *WebGeeContext) SetWriter(writer http.ResponseWriter) {
	c.writer = writer
}
func (c *WebGeeContext) GetWriter() http.ResponseWriter {
	return c.writer
}
func (c *WebGeeContext) SetStatus(code int) {
	c.statusCode = code
}
func (c *WebGeeContext) GetStatus() int {
	return c.statusCode
}
func (c *WebGeeContext) SetRequest(req *http.Request) {
	c.req = req
}
func (c *WebGeeContext) GetRequest() *http.Request {
	return c.req
}
func (c *WebGeeContext) SetPath(path string) {
	c.path = path
}
func (c *WebGeeContext) GetPath() string {
	return c.path
}
func (c *WebGeeContext) SetMethod(method string) {
	c.method = method
}
func (c *WebGeeContext) GetMethod() string {
	return c.method
}
