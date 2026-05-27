package theme

import (
	"html/template"
	"strings"
	"time"
)

// defaultFuncs returns the standard template function map available to all themes.
func defaultFuncs() template.FuncMap {
	return template.FuncMap{
		"formatDate":    formatDate,
		"formatISO":     formatISO,
		"upper":         strings.ToUpper,
		"lower":         strings.ToLower,
		"trim":          strings.TrimSpace,
		"replace":       strings.ReplaceAll,
		"startsWith":    strings.HasPrefix,
		"endsWith":      strings.HasSuffix,
		"relURL":        relURL,
		"absURL":        absURL,
		"truncate":      truncate,
		"safeHTML":      safeHTML,
		"default":       defaultVal,
		"join":          joinSep,
		"dict":          dict,
		"safeCSS":       safeCSS,
		"now":           func() time.Time { return time.Now() },
		"substr":        substr,
	}
}

func formatDate(fmt string, t time.Time) string {
	return t.Format(convertGoLayout(fmt))
}

func formatISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

// convertGoLayout converts common strftime-like patterns to Go time layout.
func convertGoLayout(s string) string {
	replacements := map[string]string{
		"%Y": "2006", "%m": "01", "%d": "02",
		"%H": "15", "%M": "04", "%S": "05",
	}
	for k, v := range replacements {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

func relURL(root, path string) string {
	if root == "" || root == "/" {
		return path
	}
	return strings.TrimRight(root, "/") + "/" + strings.TrimLeft(path, "/")
}

func absURL(baseURL, path string) string {
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func truncate(max int, s string) string {
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max]) + "…"
	}
	return s
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func safeCSS(s string) template.CSS {
	return template.CSS(s)
}

func defaultVal(def, val any) any {
	if val == nil || val == "" || val == 0 {
		return def
	}
	return val
}

func joinSep(sep string, parts []string) string {
	return strings.Join(parts, sep)
}

// dict creates a map from alternating key-value pairs.
// Usage: {{dict "key1" "val1" "key2" "val2"}}
func dict(args ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if ok {
			m[key] = args[i+1]
		}
	}
	return m
}

// substr returns a substring of s by rune index [start, end).
// If end > len(s), it is clamped. Negative indices count from end.
func substr(start, end int, s string) string {
	runes := []rune(s)
	l := len(runes)
	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}
	if start < 0 {
		start = 0
	}
	if end > l {
		end = l
	}
	if start >= end {
		return ""
	}
	return string(runes[start:end])
}
