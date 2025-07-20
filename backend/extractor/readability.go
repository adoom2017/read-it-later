package extractor

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"read-it-later/backend/model"

	"github.com/go-shiori/go-readability"
)

// Extract fetches the content from a URL and uses go-readability to parse it
// into a structured Article object.
func Extract(urlString string) (model.Article, error) {
	// Parse the URL string
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		return model.Article{}, err
	}

	// Check if this is a Zhihu article and use headless browser
	if strings.Contains(urlString, "zhihu.com") {
		browserExtractor := NewHeadlessBrowserExtractor()
		article, err := browserExtractor.ExtractWithBrowser(urlString)
		if err == nil && article.Title != "" {
			return article, nil
		}
		// If browser extraction fails, continue with regular extraction
	}

	// Check if this is a WeChat article and use headless browser
	if strings.Contains(urlString, "mp.weixin.qq.com") {
		browserExtractor := NewHeadlessBrowserExtractor()
		article, err := browserExtractor.ExtractWithBrowser(urlString)
		if err == nil && article.Title != "" {
			return article, nil
		}
		// If browser extraction fails, continue with regular extraction
	}

	// Make HTTP request with better headers
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return model.Article{}, err
	}

	// Add user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return model.Article{}, err
	}
	defer resp.Body.Close()

	// Use go-readability to parse the document
	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return model.Article{}, err
	}

	// Check if extraction was successful (content is meaningful)
	if isContentEmpty(article) {
		// Try fallback method for problematic sites
		return extractWithFallback(urlString, parsedURL)
	}

	// Create our own Article model from the parsed data
	result := model.Article{
		URL:      urlString,
		Title:    article.Title,
		Content:  article.TextContent, // Using TextContent for a cleaner reading view
		Excerpt:  article.Excerpt,
		ImageURL: ProcessImageURL(article.Image),
	}

	return result, nil
}

// isContentEmpty checks if the extracted content is meaningful
func isContentEmpty(article readability.Article) bool {
	// Check if title is empty or too short
	if len(strings.TrimSpace(article.Title)) < 3 {
		return true
	}

	// Check if content is empty or too short
	content := strings.TrimSpace(article.TextContent)
	if len(content) < 50 {
		return true
	}

	// Check if content seems to be error messages or placeholders
	lowerContent := strings.ToLower(content)
	errorKeywords := []string{
		"javascript", "请开启", "loading", "error", "404", "403",
		"access denied", "页面不存在", "内容加载中", "请稍后",
		"当前环境异常", "完成验证后即可继续访问",
	}

	for _, keyword := range errorKeywords {
		if strings.Contains(lowerContent, keyword) && len(content) < 200 {
			return true
		}
	}

	return false
}

// extractWithFallback attempts to extract content using alternative methods
func extractWithFallback(urlString string, parsedURL *url.URL) (model.Article, error) {
	// Try to get at least the title from the URL or domain
	title := getBasicTitle(parsedURL)

	result := model.Article{
		URL:      urlString,
		Title:    title,
		Content:  fmt.Sprintf("无法自动提取此页面的内容。这可能是因为该网站使用了JavaScript动态加载内容或设置了访问限制。\n\n请点击下方的 'Read Original' 链接查看原文。\n\n原文链接：%s", urlString),
		Excerpt:  "无法自动提取内容，请查看原文",
		ImageURL: "",
	}

	return result, nil
}

// getBasicTitle tries to generate a reasonable title from the URL
func getBasicTitle(parsedURL *url.URL) string {
	// Try to extract title from URL path or use domain name
	if parsedURL.Path != "" && parsedURL.Path != "/" {
		// Clean up the path to make a reasonable title
		pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		if len(pathParts) > 0 {
			lastPart := pathParts[len(pathParts)-1]
			// Remove file extensions and clean up
			if strings.Contains(lastPart, ".") {
				lastPart = strings.Split(lastPart, ".")[0]
			}
			// Replace hyphens and underscores with spaces
			lastPart = strings.ReplaceAll(lastPart, "-", " ")
			lastPart = strings.ReplaceAll(lastPart, "_", " ")
			if len(lastPart) > 5 {
				return strings.Title(lastPart)
			}
		}
	}

	// Fallback to domain name
	domain := strings.TrimPrefix(parsedURL.Host, "www.")

	return fmt.Sprintf("文章来自 %s", domain)
}
