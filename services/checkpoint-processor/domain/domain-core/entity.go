package domain

import "fmt"

type Invocation struct {
	InvocationID InvocationID
	OutputURL    OutputURL
}

func (i *Invocation) UpdateOutputURL(url OutputURL) error {
	if i.OutputURL != "" {
		return fmt.Errorf("Output URL is immutable after being set")
	}

	i.OutputURL = url

	return nil
}
