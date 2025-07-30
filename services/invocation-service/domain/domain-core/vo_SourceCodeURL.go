package domain

type SourceCodeURL string

func NewSourceCodeURL(url string) SourceCodeURL {
	return SourceCodeURL(url)
}

func (u *SourceCodeURL) String() string {
	return string(*u)
}
