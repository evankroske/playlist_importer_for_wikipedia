package unwrap

type UnwrapError struct {
	msg string
}

func (e *UnwrapError) Error() string {
	return e.msg
}

func Unwrap(t interface{}, q string) (interface{}, error) {
	return "unwrapped!", nil
}
