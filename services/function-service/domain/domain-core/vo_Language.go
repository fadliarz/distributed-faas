package domain

type Language string

func NewLanguage(language string) Language {
	return Language(language)
}

func (u *Language) String() string {
	return string(*u)
}
