package dhallconfig

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
)

func GetDhallType(config interface{}) (string, error) {
	var output bytes.Buffer

	v := reflect.ValueOf(config)
	dhallType := getDhallType(v.Type())

	cmd := exec.Command("dhall", "format")
	cmd.Stdin = strings.NewReader(dhallType)
	cmd.Stdout = &output

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.Trim(output.String(), " \n"), nil
}

func getDhallType(t reflect.Type) string {

	for {
		if t.Kind() != reflect.Ptr {
			break
		}
		t = t.Elem()
	}

	if t.PkgPath()+"."+t.Name() == "math/big.Int" {
		return "Integer"
	}

	switch t.Kind() {
	case reflect.Slice:
		return fmt.Sprintf("List %s", getDhallType(t.Elem()))
	case reflect.Struct:

		var buf bytes.Buffer

		_, _ = fmt.Fprintf(&buf, "{")

		for i := 0; i < t.NumField(); i++ {
			terminator := ", "
			if i == t.NumField()-1 {
				terminator = ""
			}

			f := t.Field(i)
			_, _ = fmt.Fprintf(&buf, "%s : %s%s", f.Name, getDhallType(f.Type), terminator)
		}
		_, _ = fmt.Fprintf(&buf, " }")
		return buf.String()
	case reflect.Bool:
		return "Bool"
	case reflect.Uint64:
		return "Natural"
	case reflect.Int64:
		return "Integer"
	case reflect.Float64:
		return "Double"
	case reflect.String:
		return "Text"
	default:
		panic(fmt.Sprintf("unsupported type %v", t))
	}

}

// LoadConfig is a convenience function, Load a dhall config while
// forwarding any errors to stderr. This is usually what you want
// for command line programs such as servers.
func LoadConfig(configExpression string, config interface{}) error {

	if configExpression == "" {
		return errors.New("Expected a config expression, got an empty string.")
	}

	var input bytes.Buffer
	var output bytes.Buffer
	var stderr bytes.Buffer

	t, err := GetDhallType(config)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(&input, "%s : %s", configExpression, t)
	cmd := exec.Command("dhall-to-json")
	cmd.Stdin = &input
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	err = json.Unmarshal(output.Bytes(), config)
	if err != nil {
		return err
	}

	return nil
}
