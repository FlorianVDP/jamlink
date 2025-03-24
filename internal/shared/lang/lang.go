package lang

import "strings"

type LangNormalizer interface {
	Normalize(raw string) string
}

type normalizer struct{}

func NewLangNormalizer() LangNormalizer {
	return &normalizer{}
}

func (n *normalizer) Normalize(raw string) string {
	if raw == "" {
		return "fr-FR"
	}

	parts := strings.Split(raw, ",")
	code := strings.TrimSpace(parts[0])

	code = formatToBCP47(code)

	if Supported[code] {
		return code
	}

	return "fr-FR"
}

func formatToBCP47(code string) string {
	segs := strings.Split(code, "-")
	if len(segs) == 2 {
		return strings.ToLower(segs[0]) + "-" + strings.ToUpper(segs[1])
	}
	return strings.ToLower(code)
}
