package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
	"wgpro/internal"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
func main() {
	fmt.Println("	启动WebGee")

	wg := internal.NewWebGeeEngine(internal.NewWebGeeRouter())

	wg.GET("/", func(c internal.IContext) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>", nil)
	})
	wg.POST("/hello", func(c internal.IContext) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.GetPath())
	})

	wg.GET("/index", func(c internal.IContext) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>", nil)
	})
	wg.Use(internal.Logger()) //添加中间件
	wg.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	wg.LoadHTMLGlob("templates/*")
	wg.Static("/assets", "./static")
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	wg.GET("/", func(c internal.IContext) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	wg.GET("/students", func(c internal.IContext) {
		c.HTML(http.StatusOK, "arr.tmpl", internal.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	wg.GET("/date", func(c internal.IContext) {
		c.HTML(http.StatusOK, "custom_func.tmpl", internal.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	wg.RUN(":9999")
}
