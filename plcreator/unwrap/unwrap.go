package unwrap

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type UnwrapError struct {
	msg string
}

func (e *UnwrapError) Error() string {
	return e.msg
}

func Unwrap(tree interface{}, q string) (interface{}, error) {
	s := new(scanner.Scanner)
	s.Init(strings.NewReader(q))
	t := s.Scan()
	for t != scanner.EOF {
		switch t {
		case '.':
			t = s.Scan()
			if t != scanner.Ident {
				return nil, &UnwrapError{
					fmt.Sprintf(
						"unwrap: Parse error: expected id, got %v",
						scanner.TokenString(t),
					),
				}
			}
			m, ok := tree.(map[string]interface{})
			if !ok {
				return nil, &UnwrapError{"unwrap: Type error"}
			}
			if tree == nil {
				return nil, &UnwrapError{"unwrap: Nil error"}
			}
			tree = m[s.TokenText()]
		case '[':
			t = s.Scan()
			i, _ := strconv.Atoi(s.TokenText())
			t = s.Scan()
			if t != ']' {
				return nil, &UnwrapError{
					fmt.Sprintf(
						"unwrap: Parse error: expected ']', got %v",
						scanner.TokenString(t),
					),
				}
			}
			a := tree.([]interface{})
			tree = a[i]
		}
		t = s.Scan()
	}
	return tree, nil
}
