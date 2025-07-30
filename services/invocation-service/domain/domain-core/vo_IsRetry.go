package domain

type IsRetry bool

func NewIsRetry(isRetry bool) IsRetry {
	return IsRetry(isRetry)
}

func (t IsRetry) Bool() bool {
	return bool(t)
}
