package envx

import (
	"os"
	"reflect"
	"strconv"

	"github.com/alvii147/lobo/errfmt"
)

// EnvTag is the struct tag used to map a struct field to an environment variable.
const EnvTag = "env"

// LoadConfig reads environment variables and populates a given config struct.
// It returns error when an environment variable is not found or properly parsed.
func LoadConfig[T any](config T) error {
	fields := reflect.TypeOf(config)
	values := reflect.ValueOf(config)

	for i := range fields.Elem().NumField() {
		field := fields.Elem().Field(i)
		value := values.Elem().Field(i)

		envKey, ok := field.Tag.Lookup(EnvTag)
		if !ok {
			continue
		}

		envValue, ok := os.LookupEnv(envKey)
		if !ok {
			return errfmt.Errorf(nil, "os.LookupEnv failed, missing environment variable %s", envKey)
		}

		switch field.Type.Kind() {
		case reflect.Bool:
			envBool, err := strconv.ParseBool(envValue)
			if err != nil {
				return errfmt.Errorf(err, "strconv.ParseBool failed to parse %s", envValue)
			}

			value.SetBool(envBool)
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			envInt, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return errfmt.Errorf(err, "strconv.ParseInt failed to parse %s", envValue)
			}

			value.SetInt(envInt)
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr:
			envInt, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				return errfmt.Errorf(err, "strconv.ParseInt failed to parse %s", envValue)
			}

			value.SetUint(envInt)
		case reflect.Float32, reflect.Float64:
			envFloat, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return errfmt.Errorf(err, "strconv.ParseFloat failed to parse %s", envValue)
			}

			value.SetFloat(envFloat)
		case reflect.String:
			value.SetString(envValue)
		default:
			return errfmt.Errorf(nil, "unsupported type %v for field %s", field.Type.Kind(), field.Name)
		}
	}

	return nil
}
