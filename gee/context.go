package gee

//包含需要的所有信息
import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}
type Context struct {
	//http原始信息
	Writer http.ResponseWriter
	Req    *http.Request
	//请求信息
	Path   string
	Method string
	//返回信息
	StatusCode int
	//动态解析的内容
	Params map[string]string

	//中间件的功能增添
	index    int
	handlers []HandlerFunc
	engine   *Engine
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

// 获取解析的内同
func (c *Context) Param(key string) string {
	value, ok := c.Params[key]
	if !ok {
		return ""
	}
	return value
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	c := Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
	return &c
}

// 获取表单数据
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取url的查询参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置返回的状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
}

// 设置响应头的格式
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 返回纯文本的字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.StatusCode = code
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回Json数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 响应二进制返回
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 响应HTML返回值
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
