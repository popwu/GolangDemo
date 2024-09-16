package main

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// 使用 sessions 中间件
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 使用 CSRF 中间件
	r.Use(csrf.Middleware(csrf.Options{
		Secret: "jsdfha0syq90ah08d",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	r.GET("/", func(c *gin.Context) {
		// 生成 CSRF 令牌
		token := csrf.GetToken(c)

		// 渲染模板
		c.HTML(200, "index.html", gin.H{
			"csrf_token": token,
		})
	})

	r.POST("/submit", func(c *gin.Context) {
		// 验证 CSRF 令牌
		// 在使用 csrf.Middleware 之后，如果 CSRF 令牌不匹配，ErrorFunc 将被调用

		// 处理表单数据
		// ...

		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	fmt.Println("Server started at http://localhost:8080")
	r.Run()
}
