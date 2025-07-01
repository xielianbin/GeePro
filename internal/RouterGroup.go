package internal

import (
	"log"
	"net/http"
	"path"
)

type RouterGroup struct {
	prefix      string        //前缀
	middlewares []HandlerFunc // 中间件处理函数
	parent      *RouterGroup  // 路由分组
	engine      *WebGeeEngine //指向共享的引擎实例
}

// 定义组来创建一个新的RouterGroup
// 所有组共享同一个Engine实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) CreateStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// 通过path.Join组合路由组前缀和相对路径得到绝对路径
	absolutePath := path.Join(group.prefix, relativePath)
	// 创建文件服务器fileServer，使用http.StripPrefix处理路径前缀
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c IContext) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.GetWriter(), c.GetRequest())
	}
}
func (group *RouterGroup) Static(relativePath string, root string) {

	handler := group.CreateStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.AddRouter(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
