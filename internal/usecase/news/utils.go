package news

import (
	"log"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func (p *NewsParser) extractFullText(item RSSItem) string {
	if item.FullText != "" {
		return item.FullText
	}
	if item.Encoded != "" {
		return item.Encoded
	}
	if item.Content != "" {
		return item.Content
	}

	return item.Description
}

func (p *NewsParser) parsePubDate(dateStr string, now time.Time) time.Time {
	if dateStr == "" {
		return now
	}

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", dateStr); err == nil {
		return t
	}

	log.Printf("не удалось распарсить дату: %s", dateStr)
	return now
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "ё", "е")
	s = regexp.MustCompile(`[^\p{L}\p{N}\s-]+`).ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func clean(s string) string {
	s = strings.TrimSpace(s)
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
}

func ptrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
