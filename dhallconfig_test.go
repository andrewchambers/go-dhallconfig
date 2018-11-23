package dhallconfig

import (
	"reflect"
	"testing"
)

type ConfigType struct {
	Foo  string
	Bar  int64
	Baz  bool
	Boo  float64
	Bag  []int64
	Nest struct{ Foo bool }
}

func TestLoadConfig(t *testing.T) {
	expectedConfig := &ConfigType{
		Foo:  "bar",
		Bar:  6,
		Baz:  true,
		Boo:  10.2,
		Bag:  []int64{1, 2},
		Nest: struct{ Foo bool }{Foo: true},
	}

	actualConfig := &ConfigType{}

	err := LoadConfig("{ Foo = \"bar\", Bar = +6, Baz = True, Boo = 10.2, Bag = [+1, +2], Nest = {Foo = True} }", actualConfig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Fatalf("expected(%#v) != %#v\n", expectedConfig, actualConfig)
	}

}
