package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main2() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: "secret123",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	r.GET("/protected", func(c *gin.Context) {
		c.String(200, csrf.GetToken(c))
	})

	r.POST("/protected", func(c *gin.Context) {
		c.String(200, "CSRF token is valid")
	})

	// req, _ := http.NewRequest("GET", "/protected", nil)
	// w := httptest.NewRecorder()

	// // 使用创建的请求
	// r.ServeHTTP(w, req)

	// // 输出响应并获取 CSRF 令牌
	// csrfToken := strings.TrimSpace(w.Body.String())
	// fmt.Println(csrfToken)

	// req, _ = http.NewRequest("POST", "/protected", nil)
	// req.Header.Set("X-CSRF-TOKEN", csrfToken)
	// w = httptest.NewRecorder()
	// r.ServeHTTP(w, req)
	// fmt.Println(w.Body.String()) // CSRF token mismatch

	// 创建一个支持 cookie 的 HTTP 客户端
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	// 发送 GET 请求并获取 CSRF 令牌
	resp, _ := client.Get("http://localhost:8080/protected")
	buf := new(strings.Builder)
	io.Copy(buf, resp.Body)
	csrfToken := strings.TrimSpace(buf.String())
	fmt.Println(csrfToken)

	// 发送 POST 请求并设置 CSRF 令牌
	req, _ := http.NewRequest("POST", "http://localhost:8080/protected", nil)
	req.Header.Set("X-CSRF-TOKEN", csrfToken)
	resp, _ = client.Do(req)
	buf = new(strings.Builder)
	io.Copy(buf, resp.Body)
	fmt.Println(buf.String()) // CSRF token is valid

}
