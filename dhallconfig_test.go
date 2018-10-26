package dhallconfig

import (
	"reflect"
	"testing"
)

type ConfigType struct {
	Foo string
	Bar int64
	Baz bool
}

func TestGetType(t *testing.T) {
	expectedConfigType := `{ Bar : Integer, Baz : Bool, Foo : Text }`
	actualConfigType, err := GetDhallType(&ConfigType{})
	if err != nil {
		t.Fatal(err)
	}
	if actualConfigType != expectedConfigType {
		t.Fatalf("expected(%s), got (%s)", expectedConfigType, actualConfigType)
	}
}

func TestLoadConfig(t *testing.T) {
	expectedConfig := &ConfigType{
		Foo: "bar",
		Bar: 6,
		Baz: true,
	}

	actualConfig := &ConfigType{}

	err := LoadConfig("{ Foo = \"bar\", Bar = 6, Baz = True }", actualConfig)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Fatalf("expected(%#v) != %#v\n", expectedConfig, actualConfig)
	}

}
