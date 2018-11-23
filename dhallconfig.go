package dhallconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"strings"

	"github.com/andrewchambers/go-extra/errors"
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
		return "", errors.Wrap(err)
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

func LoadConfig(configExpression string, config interface{}) error {

	if configExpression == "" {
		return errors.New("Expected a config expression, got an empty string.")
	}

	var input bytes.Buffer
	var output bytes.Buffer
	var stderr bytes.Buffer

	t, err := GetDhallType(config)
	if err != nil {
		return errors.Wrap(err)
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
		return errors.Wrap(err)
	}

	return nil
}
