package wrapping

import "github.com/pkg/errors"

type Error string

func (w Error) Error() string {
	return string(w)
}

func (w Error) Wrap(err error) error {
	return errors.Wrap(w, err.Error())
}
