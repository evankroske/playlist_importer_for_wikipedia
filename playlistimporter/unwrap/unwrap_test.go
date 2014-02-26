/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
