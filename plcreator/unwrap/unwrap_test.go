package unwrap

import "testing"

func TestUnwrapMap(t *testing.T) {
	var obj interface{} = map[string]interface{}{
		"foo": map[string]interface{}{ "bar" : "baz" },
	}
	res, err := Unwrap(obj, ".foo.bar")
	if err != nil {
		t.Fatal(err)
	}
	str, ok := res.(string)
	if !ok {
		t.Fatal("Wrong result type")
	} else if str != "baz" {
		t.Fatal("Wrong result value")
	}
}

func TestUnwrapMapTypeError(t *testing.T) {
	var array interface{} = make([]interface{}, 1)
	_, err := Unwrap(array, ".foo")
	if err == nil {
		t.Fatal("Expected type error")
	}

}

func TestUnwrapMapParseError(t *testing.T) {
	var object interface{}
	_, err := Unwrap(object, "..}")
	if err == nil {
		t.Fatal("Expected parse error")
	}
}

func TestUnwrapArrayIndex(t *testing.T) {
	var array interface{} = []interface{}{ []interface{}{ 1 } }
	res, err := Unwrap(array, "[0][0]")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := res.(int); !ok {
		t.Fatal("Wrong result type")
	} else if v != 1 {
		t.Fatal("Wrong result value")
	}
}

func TestUnwrapArraySlice(t *testing.T) {
	var array interface{} = []interface{}{
		[]interface{}{ 1 },
		[]interface{}{ 2 },
	}
	res, err := Unwrap(array, "[:][0]")
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := res.([]interface{}); !ok {
		t.Fatalf("Wrong result type: %v", v)
	} else if v[0] != 1 || v[1] != 2 {
		t.Fatal("Wrong result type")
	}
}
