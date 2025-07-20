package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyImage 代理图片请求，绕过防盗链
func ProxyImage(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image URL is required"})
		return
	}

	// 验证URL
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	// 安全检查：只允许代理特定域名的图片
	allowedDomains := []string{
		"mmbiz.qpic.cn",
		"wx.qpic.cn",
		"mmbiz.qlogo.cn",
	}

	allowed := false
	for _, domain := range allowedDomains {
		if strings.Contains(parsedURL.Host, domain) {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Domain not allowed for proxy"})
		return
	}

	// 创建HTTP客户端
	client := &http.Client{}
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 设置合适的请求头，模拟来自微信的请求
	if strings.Contains(parsedURL.Host, "mmbiz.qpic.cn") ||
		strings.Contains(parsedURL.Host, "wx.qpic.cn") {
		req.Header.Set("Referer", "https://mp.weixin.qq.com/")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 MicroMessenger/6.7.3.9001")
	} else {
		// 对其他域名使用通用请求头
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch image from source"})
		return
	}

	// 设置响应头
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // 默认类型
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400") // 缓存1天
	c.Header("Access-Control-Allow-Origin", "*")       // 允许跨域

	// 转发图片数据
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		// 错误已经开始写入响应，无法返回JSON错误
		return
	}
}
