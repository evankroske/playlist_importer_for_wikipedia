package unwrap

import "testing"

func TestUnwrapMap(t *testing.T) {
	var obj interface{} = map[string]int{ "foo": 1 }
	res, err := Unwrap(obj, ".foo")
	if err != nil {
		t.Fatal(err)
	}
	v, _ := res.(int)
	if v != 1 {
		t.Fatal("Wrong type")
	}
}
