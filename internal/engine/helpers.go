package engine

// Get language struct by language code/format string.
func GetLanguageInfo(language *string) *Language {
	for _, value := range SUPPORTED_LANGUAGES {
		if value.Format == *language {
			return value
		}
	}
	return nil
}
