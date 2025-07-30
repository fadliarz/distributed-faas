package domain

type OutputURL string

func NewOutputURL(url string) OutputURL {
	return OutputURL(url)
}

func (u *OutputURL) String() string {
	return string(*u)
}
