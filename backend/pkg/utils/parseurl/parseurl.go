package utils

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type UrlInfo struct {
	Title   string
	Favicon string
}

func GetTitleAndFaviconFromUrl(pageURL string) (*UrlInfo, error) {
	// Запрос к странице
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	html := string(body)

	var urlInfo UrlInfo

	// Извлекаем заголовок с помощью регулярки
	titleRe := regexp.MustCompile(`<title>(.*?)</title>`)
	titleMatch := titleRe.FindStringSubmatch(html)

	if len(titleMatch) > 1 {
		urlInfo.Title = titleMatch[1]
	}

	// Извлекаем favicon
	base, _ := url.Parse(pageURL)

	// Ищем favicon в ссылках
	faviconRe := regexp.MustCompile(`<link[^>]*rel=["'](?:shortcut\s+)?icon["'][^>]*href=["']([^"']+)["']`)
	faviconMatch := faviconRe.FindStringSubmatch(html)

	if len(faviconMatch) > 1 {
		faviconURL, err := url.Parse(faviconMatch[1])
		if err == nil {
			urlInfo.Favicon = base.ResolveReference(faviconURL).String()
		}
	} else {
		// Пробуем стандартный путь
		urlInfo.Favicon = base.Scheme + "://" + base.Host + "/favicon.ico"
	}

	return &urlInfo, nil
}
