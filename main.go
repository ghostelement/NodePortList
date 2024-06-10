package main

import (
	"NodePortList/svc/getport"

	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 创建一个带过期时间的缓存
type cacheItem struct {
	value  interface{}
	expire time.Time
}

var cache = make(map[string]cacheItem)
var mu sync.Mutex

func main() {

	r := gin.Default()

	// 读取系统环境变量
	title := os.Getenv("CLUSTER_NAME")
	if title == "" {
		title = "NodePortList"
	} else {
		title = title + "  NodePortList"
	}
	// 定义中间件
	retryLimiter := func(c *gin.Context) {
		key := c.Request.URL.Path // 使用请求路径作为缓存的key

		mu.Lock()
		item, exists := cache[key]
		now := time.Now()

		// 如果缓存存在且未过期，则直接返回缓存结果
		if exists && item.expire.After(now) {
			// 直接从缓存中读取并返回HTML响应
			c.HTML(200, "nodeport.html", gin.H{
				"title":     title,
				"NodePorts": item.value,
				"FromCache": true,
			})
			c.Abort() // 终止后续处理
		}
		mu.Unlock()

		c.Next() // 继续处理请求
	}
	// 应用中间件
	r.Use(retryLimiter)

	// 加载html模板
	r.LoadHTMLGlob("templates/html/*")

	// 配置静态web目录  第一个参数表示路由,第二个参数表示映射的路径.
	r.Static("/templates", "./templates")
	r.GET("/nodeport", func(c *gin.Context) {
		nodeports, err := getport.GetNodePort("k8s.yaml")
		// 将结果存入缓存，设置30秒过期时间
		mu.Lock()
		cache[c.Request.URL.Path] = cacheItem{value: nodeports, expire: time.Now().Add(30 * time.Second)}
		mu.Unlock()

		if err != nil {
			c.HTML(400, "error.html", gin.H{
				"Error": err.Error(),
			})
		} else {
			c.HTML(200, "nodeport.html", gin.H{
				"title":     title,
				"NodePorts": nodeports,
				"Error":     err,
				"FromCache": false,
			})
		}
	})
	r.Run(":9913") // listen and serve on 0.0.0.0:9913
}
