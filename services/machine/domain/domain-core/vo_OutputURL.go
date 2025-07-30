package domain

type OutputURL string

func NewOutputURL(url string) OutputURL {
	return OutputURL(url)
}

func (o OutputURL) String() string {
	return string(o)
}
