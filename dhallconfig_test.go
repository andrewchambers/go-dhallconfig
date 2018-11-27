package dhallconfig

import (
	"math/big"
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
	Big  *big.Int
}

func TestLoadConfig(t *testing.T) {
	expectedConfig := &ConfigType{
		Foo:  "bar",
		Bar:  6,
		Baz:  true,
		Boo:  10.2,
		Bag:  []int64{1, 2},
		Nest: struct{ Foo bool }{Foo: true},
		Big:  big.NewInt(0),
	}

	bigText := "123450000000000000000000000000000000000000000000000"
	err := expectedConfig.Big.UnmarshalText([]byte(bigText))
	if err != nil {
		t.Fatal(err)
	}

	actualConfig := &ConfigType{}

	err = LoadConfig(`{
		Foo = "bar",
		Bar = +6,
		Baz = True,
		Boo = 10.2,
		Bag = [+1, +2],
		Nest = {Foo = True},
		Big = +123450000000000000000000000000000000000000000000000
		}`, actualConfig)
	if err != nil {
		t.Fatal(err)
	}

	if actualConfig.Foo != expectedConfig.Foo {
		t.Fatal("bad value.")
	}

	if actualConfig.Bar != expectedConfig.Bar {
		t.Fatal("bad value.")
	}

	if actualConfig.Baz != expectedConfig.Baz {
		t.Fatal("bad value.")
	}

	if actualConfig.Boo != expectedConfig.Boo {
		t.Fatal("bad value.")
	}

	if !reflect.DeepEqual(actualConfig.Bag, expectedConfig.Bag) {
		t.Fatalf("bad value.")
	}

	if !reflect.DeepEqual(actualConfig.Nest, expectedConfig.Nest) {
		t.Fatal("bad value.")
	}

	if !reflect.DeepEqual(actualConfig.Nest, expectedConfig.Nest) {
		t.Fatal("bad value.")
	}

	if actualConfig.Big.Text(10) != expectedConfig.Big.Text(10) {
		t.Fatalf("bad value. %v , %v", actualConfig.Big.Text(10), expectedConfig.Big.Text(10))
	}

}
