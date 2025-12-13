package parseurl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UrlInfo interface {
	GetTitle() string
	GetFaviconPath() string
	DownloadFavicon(saveDir string, userID, linkID int) (string, error)
}

type urlInfo struct {
	url     string
	title   string
	favicon string
}

func New(rawURL string) UrlInfo {
	out := &urlInfo{url: normalizeURL(rawURL)}
	
	// Загружаем данные синхронно
	_ = out.loadData()
	
	return out
}

// normalizeURL добавляет схему протокола если отсутствует
func normalizeURL(rawURL string) string {
	// Удаляем пробелы
	rawURL = strings.TrimSpace(rawURL)
	
	// Если URL уже содержит схему, возвращаем как есть
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		return rawURL
	}
	
	// Добавляем https:// по умолчанию
	return "https://" + rawURL
}

func (u *urlInfo) GetTitle() string {
	return u.title
}

func (u *urlInfo) GetFaviconPath() string {
	return u.favicon
}

func (u *urlInfo) loadData() error {
	// Проверяем валидность URL
	parsedURL, err := url.Parse(u.url)
	if err != nil {
		return fmt.Errorf("неверный URL: %w", err)
	}
	
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
		u.url = parsedURL.String()
	}

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Запрос к странице
	resp, err := client.Get(u.url)
	if err != nil {
		return fmt.Errorf("ошибка HTTP запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP статус: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %w", err)
	}
	html := string(body)

	// Извлекаем заголовок с помощью регулярки
	titleRe := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	titleMatch := titleRe.FindStringSubmatch(html)

	if len(titleMatch) > 1 {
		u.title = strings.TrimSpace(titleMatch[1])
	}

	// Находим фавиконку
	u.favicon = findFavicon(html, parsedURL)

	return nil
}

// Функция для скачивания и сохранения favicon
// ВОЗВРАЩАЕТ ЛОКАЛЬНЫЙ ПУТЬ К ФАЙЛУ НА ДИСКЕ
func (u *urlInfo) DownloadFavicon(saveDir string, userID, linkID int) (string, error) {
	if u.favicon == "" {
		// Если favicon не найден, возвращаем пустую строку
		return "", nil
	}

	// Создаем путь для сохранения: saveDir/userID/
	userDir := filepath.Join(saveDir, strconv.Itoa(userID))

	// Создаем директорию пользователя, если не существует
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию пользователя: %w", err)
	}

	// Определяем расширение файла из URL favicon
	faviconURL, err := url.Parse(u.favicon)
	if err != nil {
		return "", fmt.Errorf("неверный URL favicon: %w", err)
	}

	// Получаем расширение из URL
	ext := strings.ToLower(filepath.Ext(faviconURL.Path))

	// Если расширение не найдено в пути, пробуем определить из Content-Type
	if ext == "" {
		ext = u.detectExtensionFromURL(u.favicon)
	}

	// Убираем возможные параметры из расширения (например, .ico?v=2)
	if idx := strings.Index(ext, "?"); idx != -1 {
		ext = ext[:idx]
	}
	if idx := strings.Index(ext, "#"); idx != -1 {
		ext = ext[:idx]
	}

	// Проверяем, что расширение корректное
	validExts := map[string]bool{
		".ico": true, ".png": true, ".jpg": true,
		".jpeg": true, ".gif": true, ".svg": true, ".webp": true,
	}

	if !validExts[ext] {
		// Если расширение не валидное, используем .ico
		ext = ".ico"
	}

	// Формируем имя файла: linkID + расширение
	filename := strconv.Itoa(linkID) + ext
	filePath := filepath.Join(userDir, filename)

	// УДАЛЯЕМ СУЩЕСТВУЮЩИЙ ФАЙЛ, ЕСЛИ ОН ЕСТЬ
	// Это нужно чтобы обновить иконку если логотип изменился
	if _, err := os.Stat(filePath); err == nil {
		// Файл существует - удаляем его
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Предупреждение: не удалось удалить старый файл %s: %v\n", filePath, err)
			// Продолжаем скачивание, возможно перезапишем файл
		}
	}

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Скачиваем favicon
	resp, err := client.Get(u.favicon)
	if err != nil {
		return "", fmt.Errorf("ошибка при загрузке favicon: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP ошибка: %s", resp.Status)
	}

	// Проверяем Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("неверный Content-Type: %s", contentType)
	}

	// Создаем файл (os.Create всегда создает новый или перезаписывает существующий)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	// Копируем данные
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		// Удаляем файл в случае ошибки
		os.Remove(filePath)
		return "", fmt.Errorf("ошибка при сохранении файла: %w", err)
	}

	// Проверяем, что файл не пустой
	fileInfo, err := file.Stat()
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("ошибка получения информации о файле: %w", err)
	}
	
	if fileInfo.Size() == 0 {
		os.Remove(filePath)
		return "", fmt.Errorf("скачанный файл пуст")
	}

	// ВОЗВРАЩАЕМ ЛОКАЛЬНЫЙ ПУТЬ К ФАЙЛУ
	// Этот путь нужно сохранить в базу в поле FaviconURL
	return filePath, nil
}

// detectExtensionFromURL определяет расширение файла по заголовкам HTTP
func (u *urlInfo) detectExtensionFromURL(faviconURL string) string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	req, err := http.NewRequest("HEAD", faviconURL, nil)
	if err != nil {
		return ".ico"
	}
	
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		switch {
		case strings.Contains(contentType, "image/x-icon"),
			strings.Contains(contentType, "image/vnd.microsoft.icon"):
			return ".ico"
		case strings.Contains(contentType, "image/png"):
			return ".png"
		case strings.Contains(contentType, "image/jpeg"):
			return ".jpg"
		case strings.Contains(contentType, "image/gif"):
			return ".gif"
		case strings.Contains(contentType, "image/svg+xml"):
			return ".svg"
		case strings.Contains(contentType, "image/webp"):
			return ".webp"
		default:
			return ".ico"
		}
	}
	
	return ".ico"
}

// findFavicon поиск favicon
func findFavicon(html string, base *url.URL) string {
	// Список возможных путей к favicon
	possiblePaths := []string{
		// Стандартные пути
		"/favicon.ico",
		"/favicon.png",
		"/favicon.jpg",
		"/favicon.jpeg",
		"/favicon.gif",
		"/favicon.webp",
		"/apple-touch-icon.png",
		"/apple-touch-icon-precomposed.png",
		"/apple-touch-icon-180x180.png",
		"/apple-touch-icon-120x120.png",
		"/apple-touch-icon-76x76.png",
		"/apple-touch-icon-60x60.png",
	}

	// Ищем в HTML тегах <link> (регистронезависимо)
	htmlLower := strings.ToLower(html)
	
	// Паттерны для поиска
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`<link[^>]*rel=["'](?:shortcut\s+)?(?:icon|apple-touch-icon)["'][^>]*href=["']([^"']+)["']`),
		regexp.MustCompile(`<link[^>]*href=["']([^"']+)["'][^>]*rel=["'](?:shortcut\s+)?(?:icon|apple-touch-icon)["']`),
		regexp.MustCompile(`<meta[^>]*property=["']og:image["'][^>]*content=["']([^"']+)["']`),
		regexp.MustCompile(`<meta[^>]*name=["']twitter:image["'][^>]*content=["']([^"']+)["']`),
	}

	// Ищем в HTML
	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(htmlLower, -1)
		for _, match := range matches {
			if len(match) > 1 && match[1] != "" {
				faviconURL, err := url.Parse(match[1])
				if err == nil {
					fullURL := base.ResolveReference(faviconURL).String()
					if checkFaviconExists(fullURL) {
						return fullURL
					}
				}
			}
		}
	}

	// Проверяем стандартные пути
	for _, path := range possiblePaths {
		faviconURL := base.Scheme + "://" + base.Host + path
		if checkFaviconExists(faviconURL) {
			return faviconURL
		}
	}

	// Пробуем /favicon.ico как последний вариант
	defaultFavicon := base.Scheme + "://" + base.Host + "/favicon.ico"
	if checkFaviconExists(defaultFavicon) {
		return defaultFavicon
	}

	return ""
}

// checkFaviconExists проверка существования favicon
func checkFaviconExists(faviconURL string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	req, err := http.NewRequest("HEAD", faviconURL, nil)
	if err != nil {
		return false
	}
	
	// Добавляем User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}