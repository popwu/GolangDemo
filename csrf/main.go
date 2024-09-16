package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"net/http/httptest"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
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

	// 创建一个测试服务器
	ts := httptest.NewServer(r)
	defer ts.Close()

	// 创建一个支持 cookie 的 HTTP 客户端
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	fmt.Println("===", ts.URL)

	// 发送 GET 请求并获取 CSRF 令牌
	resp, _ := client.Get(ts.URL + "/protected")
	buf := new(strings.Builder)
	io.Copy(buf, resp.Body)
	csrfToken := strings.TrimSpace(buf.String())
	fmt.Println(csrfToken)

	// 发送 POST 请求并设置 CSRF 令牌
	req, _ := http.NewRequest("POST", ts.URL+"/protected", nil)
	req.Header.Set("X-CSRF-TOKEN", csrfToken)
	resp, _ = client.Do(req)
	buf = new(strings.Builder)
	io.Copy(buf, resp.Body)
	fmt.Println(buf.String()) // CSRF token is valid
}
