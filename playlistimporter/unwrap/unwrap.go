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
	s.Mode &= ^uint(scanner.ScanFloats)
	s.Mode |= scanner.ScanInts
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
			if tree == nil {
				return nil, &UnwrapError{"unwrap: Nil error"}
			}
			m, ok := tree.(map[string]interface{})
			if !ok {
				return nil, &UnwrapError{
					fmt.Sprintf(
						"unwrap: Type error: expected map[string]interface{}; got %T",
						tree,
					),
				}
			}
			tree = m[s.TokenText()]
		case '[':
			t = s.Scan()
			switch t {
			case scanner.Int:
				i, _ := strconv.Atoi(s.TokenText())
				a := tree.([]interface{})
				tree = a[i]
				_, err := expect(s, ']')
				if err != nil {
					return nil, err
				}
			case ':':
				a := tree.([]interface{})
				b := make([]interface{}, len(a))
				_, err := expect(s, ']')
				if err != nil {
					return nil, err
				}
				for i, v := range a {
					res, err := Unwrap(v, q[s.Offset:])
					if err != nil {
						return nil, err
					}
					b[i] = res
				}
				return b, nil
			default:
				return nil, &UnwrapError{
					fmt.Sprintf(
						"unwrap: Parse error: expected int or ':', got %v",
						scanner.TokenString(t),
					),
				}
			}
		}
		t = s.Scan()
	}
	return tree, nil
}

func expect(s *scanner.Scanner, tok rune) (t rune, err error) {
	t = s.Scan()
	if t != ']' {
		return -1, &UnwrapError{
			fmt.Sprintf(
				"unwrap: Parse error: expected ']', got %v",
				scanner.TokenString(t),
			),
		}
	}
	return
}
