package domain

type SourceCodeURL string

func (u *SourceCodeURL) String() string {
	return string(*u)
}

func NewSourceCodeURL(url string) (SourceCodeURL, error) {
	return SourceCodeURL(url), nil
}
