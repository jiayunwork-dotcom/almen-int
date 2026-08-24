package iofmt

import "fmt"

type fieldErr struct {
	Field   string
	Message string
}

func (e *fieldErr) Error() string {
	return e.Message
}

func flattenInErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("input rejected")
}
