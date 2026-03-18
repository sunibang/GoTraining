package strings

import "errors"

type Strings struct {
	prefix string
}

// AddPrefix adds the prefix stored in the Strings struct to the provided string
// and returns the result along with an error if the provided string is empty
func (s *Strings) AddPrefix(str string) (string, error) {
	if str == "" {
		return "", errors.New("empty string")
	}

	return s.prefix + str, nil
}
