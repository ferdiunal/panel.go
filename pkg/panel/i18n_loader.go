package panel

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	fiberi18n "github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type mergedLocaleLoader struct {
	userRootPath string
}

func newMergedLocaleLoader(userRootPath string) fiberi18n.Loader {
	if userRootPath == "" {
		userRootPath = "./locales"
	}
	return &mergedLocaleLoader{userRootPath: userRootPath}
}

func (l *mergedLocaleLoader) LoadMessage(requestPath string) ([]byte, error) {
	fileName := filepath.Base(requestPath)
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
	if ext == "" {
		ext = "yaml"
	}

	lang := strings.TrimSuffix(fileName, "."+ext)
	if lang == "" {
		return nil, fmt.Errorf("invalid locale file path: %s", requestPath)
	}

	merged := make(map[string]interface{})
	loaded := false

	for _, code := range localeLanguageCodes(lang) {
		for _, candidateExt := range preferredEmbeddedLocaleExtensions(ext) {
			embeddedPath := path.Join("locales", code+"."+candidateExt)
			if data, err := fs.ReadFile(embeddedLocalesFS, embeddedPath); err == nil {
				parsed, parseErr := parseLocaleBytes(data, candidateExt)
				if parseErr != nil {
					return nil, parseErr
				}
				deepMergeLocaleMaps(merged, parsed)
				loaded = true
			}
		}
	}

	for _, code := range localeLanguageCodes(lang) {
		for _, candidateExt := range preferredUserLocaleExtensions(ext) {
			userPath := filepath.Join(l.userRootPath, code+"."+candidateExt)
			if data, err := os.ReadFile(userPath); err == nil {
				parsed, parseErr := parseLocaleBytes(data, candidateExt)
				if parseErr != nil {
					return nil, parseErr
				}
				deepMergeLocaleMaps(merged, parsed)
				loaded = true
			}
		}
	}

	if !loaded {
		return nil, os.ErrNotExist
	}

	return marshalLocaleBytes(merged, ext)
}

func preferredEmbeddedLocaleExtensions(primary string) []string {
	extensions := []string{primary}

	// Embedded locale fallback: SDK varsayılan dosyaları YAML olarak tutulur.
	for _, ext := range []string{"yaml", "yml"} {
		if ext != primary {
			extensions = append(extensions, ext)
		}
	}

	return extensions
}

func preferredUserLocaleExtensions(primary string) []string {
	return []string{primary}
}

func localeLanguageCodes(lang string) []string {
	normalized := strings.TrimSpace(strings.ReplaceAll(lang, "_", "-"))
	if normalized == "" {
		return nil
	}

	normalized = strings.ToLower(normalized)
	base := normalized
	if idx := strings.Index(normalized, "-"); idx > 0 {
		base = normalized[:idx]
	}

	codes := []string{base}
	if normalized != base {
		codes = append(codes, normalized)
	}

	return codes
}

func parseLocaleBytes(data []byte, format string) (map[string]interface{}, error) {
	var parsed map[string]interface{}

	switch strings.ToLower(format) {
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
	case "json":
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
	case "toml":
		if err := toml.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported locale format: %s", format)
	}

	return normalizeLocaleMap(parsed), nil
}

func marshalLocaleBytes(data map[string]interface{}, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "yaml", "yml":
		return yaml.Marshal(data)
	case "json":
		return json.Marshal(data)
	case "toml":
		return toml.Marshal(data)
	default:
		return nil, fmt.Errorf("unsupported locale format: %s", format)
	}
}

func deepMergeLocaleMaps(dest, src map[string]interface{}) {
	for key, srcValue := range src {
		if destValue, exists := dest[key]; exists {
			destMap, destOk := destValue.(map[string]interface{})
			srcMap, srcOk := srcValue.(map[string]interface{})
			if destOk && srcOk {
				deepMergeLocaleMaps(destMap, srcMap)
				dest[key] = destMap
				continue
			}
		}
		dest[key] = srcValue
	}
}

func normalizeLocaleMap(input map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(input))
	for key, value := range input {
		normalized[key] = normalizeLocaleValue(value)
	}
	return normalized
}

func normalizeLocaleValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return normalizeLocaleMap(v)
	case map[interface{}]interface{}:
		normalized := make(map[string]interface{}, len(v))
		for mk, mv := range v {
			normalized[fmt.Sprint(mk)] = normalizeLocaleValue(mv)
		}
		return normalized
	case []interface{}:
		result := make([]interface{}, 0, len(v))
		for _, item := range v {
			result = append(result, normalizeLocaleValue(item))
		}
		return result
	default:
		return value
	}
}
