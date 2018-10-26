package dhallconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"

	"github.com/andrewchambers/go-extra/errors"
)

func GetDhallType(config interface{}) string {
	var b bytes.Buffer
	getDhallType(&b, config)
	return b.String()
}

func getDhallType(buf *bytes.Buffer, config interface{}) {

	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	_, _ = fmt.Fprintf(buf, "{ ")

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		terminator := ", "
		if i == t.NumField()-1 {
			terminator = ""
		}

		f := t.Field(i)
		k := f.Type.Kind()
		switch k {
		case reflect.Bool:
			_, _ = fmt.Fprintf(buf, "%s : Bool%s", f.Name, terminator)
		case reflect.Uint64:
			_, _ = fmt.Fprintf(buf, "%s : Natural%s", f.Name, terminator)
		case reflect.Int64:
			_, _ = fmt.Fprintf(buf, "%s : Integer%s", f.Name, terminator)
		case reflect.String:
			_, _ = fmt.Fprintf(buf, "%s : Text%s", f.Name, terminator)
		default:
			panic(fmt.Sprintf("unsupported kind %v", k))
		}
	}

	_, _ = fmt.Fprintf(buf, " }")
}

func LoadConfig(configExpression string, config interface{}) error {

	var input bytes.Buffer
	var output bytes.Buffer
	var stderr bytes.Buffer

	_, _ = fmt.Fprintf(&input, "%s : %s", configExpression, GetDhallType(config))
	cmd := exec.Command("dhall-to-json")
	cmd.Stdin = &input
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	err = json.Unmarshal(output.Bytes(), config)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
