package cnf

import "strings"

// normalizeEnvVarKey returns normalized string in format `word.word2.word3` - all lower case and dot separated - or
// empty string if can't be normalized.
func normalizeEnvVarKey(s string) string {
	// Although, technically valid env variable name may start with underscore, we're going to ignore such variables.
	// We also do ignore names ending with underscore, as well as any name containing multiple underscores in a row.
	// Basically, we need a format like `WORD_word2_WORD3` case insensitive.
	///
	// See for allowed chars:
	// https://stackoverflow.com/questions/2821043/allowed-characters-in-linux-environment-variable-names
	if strings.HasPrefix(s, "_") || strings.HasSuffix(s, "_") || strings.Contains(s, "__") {
		return ""
	}
	sLower := strings.ToLower(s)
	split := strings.ReplaceAll(sLower, "_", ".")
	return split
}

func normalizeCmdLineArgKey(arg string) string {
	if strings.HasSuffix(arg, ".") || strings.HasSuffix(arg, "-") || strings.Contains(arg, "..") {
		return ""
	}
	words := strings.Split(arg, ".")
	var normalized []string
	for _, w := range words {
		normalized = append(normalized, normalizeWord(w))
	}
	return strings.Join(normalized, ".")
}

func normalizeWord(s string) string {
	var words []string
	from := 0
	for i := 1; i < len(s); i++ {
		ch := s[i]
		if ch >= 'A' && ch <= 'Z' {
			words = append(words, s[from:i])
			from = i
		}
	}
	words = append(words, s[from:])
	joined := strings.Join(words, "-")
	return strings.ToLower(joined)
}
