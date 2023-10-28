package engine

func GetLanguageInfo(language *string) *Language {
	for _, value := range SUPPORTED_LANGUAGES {
		if value.Format == *language {
			return value
		}
	}
	return nil
}
