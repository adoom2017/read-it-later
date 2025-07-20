package extractor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"read-it-later/backend/model"

	"github.com/chromedp/chromedp"
)

// HeadlessBrowserExtractor uses Chrome headless browser to extract content
type HeadlessBrowserExtractor struct {
	timeout time.Duration
}

// NewHeadlessBrowserExtractor creates a new headless browser extractor
func NewHeadlessBrowserExtractor() *HeadlessBrowserExtractor {
	return &HeadlessBrowserExtractor{
		timeout: 30 * time.Second,
	}
}

// ExtractWithBrowser extracts content using headless Chrome
func (hbe *HeadlessBrowserExtractor) ExtractWithBrowser(urlString string) (model.Article, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), hbe.timeout)
	defer cancel()

	// Create Chrome options - try to find Chrome first
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	// Create allocator context
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create browser context
	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var title, content, description, imageURL string

	// Run browser tasks
	err := chromedp.Run(browserCtx,
		// Navigate to the page
		chromedp.Navigate(urlString),

		// Wait for page to load
		chromedp.Sleep(5*time.Second),

		// Wait for content to be available (try multiple selectors)
		hbe.waitForContent(),

		// Extract title
		chromedp.Evaluate(`document.title || ''`, &title),

		// Extract meta description
		chromedp.Evaluate(`
			const metaDesc = document.querySelector('meta[name="description"]') || 
							document.querySelector('meta[property="og:description"]');
			metaDesc ? metaDesc.getAttribute('content') || '' : '';
		`, &description),

		// Extract meta image
		chromedp.Evaluate(`
			const metaImg = document.querySelector('meta[property="og:image"]') || 
							document.querySelector('meta[name="twitter:image"]');
			metaImg ? metaImg.getAttribute('content') || '' : '';
		`, &imageURL),

		// Extract main content for different sites
		hbe.extractMainContent(urlString, &content),
	)

	if err != nil {
		// If browser fails, return a meaningful error
		return model.Article{}, fmt.Errorf("headless browser not available or failed: %v", err)
	}

	// Clean and process extracted data
	article := model.Article{
		URL:      urlString,
		Title:    hbe.cleanTitle(title),
		Content:  hbe.cleanContent(content),
		Excerpt:  hbe.createExcerpt(description, content),
		ImageURL: ProcessImageURL(imageURL),
	}

	// Ensure we have meaningful content
	if len(strings.TrimSpace(article.Content)) < 50 {
		article.Content = fmt.Sprintf("通过浏览器访问成功，但提取的内容较少。\n\n页面标题：%s\n描述：%s\n\n这可能是一个需要特殊处理的页面类型。请点击查看原文获取完整内容。", article.Title, description)
	}

	return article, nil
} // waitForContent waits for content to be loaded
func (hbe *HeadlessBrowserExtractor) waitForContent() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Try to wait for different content indicators
		selectors := []string{
			// WeChat selectors
			".rich_media_content",
			"#js_content",
			".rich_media_title",

			// Zhihu selectors
			".Post-RichTextContainer",
			".RichText",
			".Post-content",
			".ContentItem-title",

			// Generic selectors
			"article",
			"main",
			".content",
			".article-content",
			"h1",
		}

		for _, selector := range selectors {
			err := chromedp.WaitVisible(selector, chromedp.ByQuery).Do(ctx)
			if err == nil {
				return nil // Found content, continue
			}
		}

		// If no specific content found, just wait a bit more
		return chromedp.Sleep(2 * time.Second).Do(ctx)
	})
}

// extractMainContent extracts main content based on the website
func (hbe *HeadlessBrowserExtractor) extractMainContent(urlString string, content *string) chromedp.Action {
	if strings.Contains(urlString, "mp.weixin.qq.com") {
		return hbe.extractWechatContent(content)
	} else if strings.Contains(urlString, "zhuanlan.zhihu.com") {
		return hbe.extractZhihuContent(content)
	} else {
		return hbe.extractGenericContent(content)
	}
}

// extractWechatContent extracts WeChat article content
func (hbe *HeadlessBrowserExtractor) extractWechatContent(content *string) chromedp.Action {
	return chromedp.Evaluate(`
		function extractWechatContent() {
			// Try multiple selectors for WeChat content
			const selectors = [
				'.rich_media_content',
				'#js_content',
				'.rich_media_area_primary .rich_media_content',
				'[data-role="main"]'
			];
			
			for (const selector of selectors) {
				const element = document.querySelector(selector);
				if (element) {
					// Clean up the content
					const clone = element.cloneNode(true);
					
					// Remove unwanted elements
					const unwanted = clone.querySelectorAll('script, style, .rich_media_tool, .rich_media_meta, [data-role="bottom"]');
					unwanted.forEach(el => el.remove());
					
					const text = clone.innerText || clone.textContent || '';
					if (text.trim().length > 100) {
						return text.trim();
					}
				}
			}
			
			// Fallback: try to get any meaningful text
			const body = document.querySelector('body');
			if (body) {
				const text = body.innerText || body.textContent || '';
				return text.trim();
			}
			
			return '';
		}
		
		extractWechatContent();
	`, content)
}

// extractZhihuContent extracts Zhihu article content
func (hbe *HeadlessBrowserExtractor) extractZhihuContent(content *string) chromedp.Action {
	return chromedp.Evaluate(`
		function extractZhihuContent() {
			// Try multiple selectors for Zhihu content
			const selectors = [
				'.Post-RichTextContainer',
				'.RichText',
				'.Post-content',
				'.ArticleItem-content',
				'[data-testid="article-content"]'
			];
			
			for (const selector of selectors) {
				const element = document.querySelector(selector);
				if (element) {
					// Clean up the content
					const clone = element.cloneNode(true);
					
					// Remove unwanted elements
					const unwanted = clone.querySelectorAll('script, style, .Post-NormalMain, .ContentItem-actions');
					unwanted.forEach(el => el.remove());
					
					const text = clone.innerText || clone.textContent || '';
					if (text.trim().length > 100) {
						return text.trim();
					}
				}
			}
			
			// Fallback: try article tag
			const article = document.querySelector('article');
			if (article) {
				const text = article.innerText || article.textContent || '';
				if (text.trim().length > 100) {
					return text.trim();
				}
			}
			
			return '';
		}
		
		extractZhihuContent();
	`, content)
}

// extractGenericContent extracts content using generic selectors
func (hbe *HeadlessBrowserExtractor) extractGenericContent(content *string) chromedp.Action {
	return chromedp.Evaluate(`
		function extractGenericContent() {
			// Try common content selectors
			const selectors = [
				'article',
				'main',
				'.content',
				'.article-content',
				'.post-content',
				'.entry-content',
				'[role="main"]'
			];
			
			for (const selector of selectors) {
				const element = document.querySelector(selector);
				if (element) {
					const text = element.innerText || element.textContent || '';
					if (text.trim().length > 100) {
						return text.trim();
					}
				}
			}
			
			// Fallback: get body text
			const body = document.querySelector('body');
			if (body) {
				// Remove script and style tags
				const clone = body.cloneNode(true);
				const unwanted = clone.querySelectorAll('script, style, nav, header, footer, aside, .sidebar, .navigation');
				unwanted.forEach(el => el.remove());
				
				const text = clone.innerText || clone.textContent || '';
				return text.trim();
			}
			
			return '';
		}
		
		extractGenericContent();
	`, content)
}

// extractFallbackContent extracts content from full HTML as fallback
func (hbe *HeadlessBrowserExtractor) extractFallbackContent(htmlContent string) string {
	// Simple text extraction from HTML
	// Remove script and style tags
	content := htmlContent

	// Basic cleanup - this could be enhanced with proper HTML parsing
	scriptStart := strings.Index(content, "<script")
	for scriptStart != -1 {
		scriptEnd := strings.Index(content[scriptStart:], "</script>")
		if scriptEnd != -1 {
			content = content[:scriptStart] + content[scriptStart+scriptEnd+9:]
		} else {
			break
		}
		scriptStart = strings.Index(content, "<script")
	}

	styleStart := strings.Index(content, "<style")
	for styleStart != -1 {
		styleEnd := strings.Index(content[styleStart:], "</style>")
		if styleEnd != -1 {
			content = content[:styleStart] + content[styleStart+styleEnd+8:]
		} else {
			break
		}
		styleStart = strings.Index(content, "<style")
	}

	// Extract text between common content tags
	if start := strings.Index(content, "<body"); start != -1 {
		if end := strings.Index(content[start:], "</body>"); end != -1 {
			bodyContent := content[start : start+end]
			// This is very basic - in a real implementation you'd use proper HTML parsing
			return bodyContent
		}
	}

	return ""
}

// cleanTitle cleans the extracted title
func (hbe *HeadlessBrowserExtractor) cleanTitle(title string) string {
	title = strings.TrimSpace(title)

	// Remove common suffixes
	suffixes := []string{
		" - 知乎",
		" - 微信公众号",
		" - WeChat",
		" - 公众号",
	}

	for _, suffix := range suffixes {
		if strings.HasSuffix(title, suffix) {
			title = strings.TrimSuffix(title, suffix)
			break
		}
	}

	return title
}

// cleanContent cleans the extracted content
func (hbe *HeadlessBrowserExtractor) cleanContent(content string) string {
	if content == "" {
		return ""
	}

	// Basic cleanup
	content = strings.TrimSpace(content)

	// Remove excessive whitespace
	lines := strings.Split(content, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// createExcerpt creates an excerpt from description or content
func (hbe *HeadlessBrowserExtractor) createExcerpt(description, content string) string {
	if description != "" && len(description) > 10 {
		if len(description) > 200 {
			return description[:200] + "..."
		}
		return description
	}

	if content != "" {
		// Create excerpt from content
		content = strings.ReplaceAll(content, "\n", " ")
		words := strings.Fields(content)

		if len(words) > 30 {
			excerpt := strings.Join(words[:30], " ")
			return excerpt + "..."
		}
		return content
	}

	return "使用浏览器提取的内容"
}
