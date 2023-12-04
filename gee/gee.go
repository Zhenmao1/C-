package gee

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

// engine结构体，保存gee框架的状态，也就是路由表
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// 分组控制的结构体
type RouterGroup struct {
	prefix      string
	middleWares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine

	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

type HandlerFunc func(c *Context)

// 构建一个框架，返回其地址
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (rg *RouterGroup) NewGroup(prefix string) *RouterGroup {
	engine := rg.engine
	new := &RouterGroup{
		prefix: rg.prefix + prefix,
		parent: rg,
		engine: engine,
	}
	engine.groups = append(engine.groups, new)
	return new
}

func (eg *Engine) SetFuncMap(funcMap template.FuncMap) {
	eg.funcMap = funcMap
}

func (eg *Engine) LoadHTMLGlob(pattern string) {
	eg.htmlTemplates = template.Must(template.New("").Funcs(eg.funcMap).ParseGlob(pattern))
}

func (rg *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := rg.prefix + relativePath
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (rg *RouterGroup) Static(relativePath string, root string) {
	handler := rg.createStaticHandler(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "/*filepath")
	rg.GET(urlPath, handler)
}

func (rg *RouterGroup) Use(handler HandlerFunc) {
	rg.middleWares = append(rg.middleWares, handler)
}

// 参数请求方法、路径和处理函数。
// 构造一个唯一的键，用于在路由映射表中存储这个路由。
// 将这个键和处理函数的映射关系添加到路由映射表中。
func (rg *RouterGroup) AddRouter(method string, comp string, handler HandlerFunc) {
	path := rg.prefix + comp
	rg.engine.router.AddRouter(method, path, handler)
}

func (rg *RouterGroup) GET(path string, handler HandlerFunc) {
	rg.AddRouter("GET", path, handler)
}
func (rg *RouterGroup) POST(path string, handler HandlerFunc) {
	rg.AddRouter("POST", path, handler)
}

func (eg *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, eg)
}
func (eg *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range eg.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleWares...)
		} else {
			// 调试输出
			fmt.Printf("req.URL.Path: %s\n", req.URL.Path)
			fmt.Printf("group.prefix: %s\n", group.prefix)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = eg
	eg.router.Handle(c)
}
