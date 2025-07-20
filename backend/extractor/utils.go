package extractor

import (
	"net/url"
	"strings"
)

// ProcessImageURL processes image URLs to handle anti-hotlinking
func ProcessImageURL(imageURL string) string {
	if imageURL == "" {
		return ""
	}

	// Handle WeChat images (mmbiz.qpic.cn) that have anti-hotlinking protection
	if strings.Contains(imageURL, "mmbiz.qpic.cn") {
		// Convert to proxy URL to bypass anti-hotlinking
		return "/api/proxy/image?url=" + url.QueryEscape(imageURL)
	}

	// Handle other potential anti-hotlinking domains
	antiHotlinkDomains := []string{
		"wx.qpic.cn",
		"mmbiz.qlogo.cn",
	}

	for _, domain := range antiHotlinkDomains {
		if strings.Contains(imageURL, domain) {
			return "/api/proxy/image?url=" + url.QueryEscape(imageURL)
		}
	}

	// For other images, return as-is
	return imageURL
}
