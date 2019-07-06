package conf

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type Resolver interface {
	Conf() map[string]interface{}
	Warns() []string
}

type genericLoader interface {
	name() string
	load()
	conf() map[string]interface{}
	warns() []string
}

// Configure injects configuration. Root must be a reference.
func Configure(root interface{}, r Resolver) error {
	// If root is nil we just ignore it
	if root == nil {
		return nil
	}
	// Must be a pointer
	v := reflect.ValueOf(root)
	if v.Kind() != reflect.Ptr {
		return errors.Errorf("root is not a pointer")
	}
	err := traverse(v, r.Conf())
	if err != nil {
		return errors.Wrap(err, "failed to inject configuration")
	}
	return nil
}

func traverse(v reflect.Value, conf map[string]interface{}) error {
	if v.IsNil() {
		panic("AAAAGRH!!!! IT IS NIL!!!!!")
		// TODO - make an obj so we can proceed
	}

	vk := v.Kind()
	if vk != reflect.Ptr && vk != reflect.Interface {
		return errors.Errorf("pointer or interface expected, got %v", vk)
	}

	var ve reflect.Value
	ve = v.Elem()

	// Must be a pointer to a struct
	vek := ve.Kind()
	if vek != reflect.Struct {
		return errors.Errorf("struct expected, got %v", vek)
	}

	// Traverse
	for i := 0; i < ve.NumField(); i++ {
		f := ve.Field(i)
		kind := f.Type().Kind()
		if kind == reflect.Struct {
			// Struct, so go deeper
			arg := f.Addr()
			err := traverse(arg, conf)
			if err != nil {
				return errors.Wrap(err, "failed to traverse conf struct")
			}
		} else {
			// Non-struct field, so try to inject
			tag := reflect.Indirect(ve).Type().Field(i).Tag
			err := inject(f, tag, conf)
			if err != nil {
				return errors.Wrap(err, "failed to inject conf value")
			}
		}
	}
	return nil
}

func inject(obj reflect.Value, tag reflect.StructTag, conf map[string]interface{}) error {
	// Only fields with proper tag allowed
	lookup, specified := tag.Lookup("conf")
	if !specified {
		// No `conf` tag, it's fine
		return nil
	}
	// Must be settable
	if !obj.CanSet() {
		return errors.New("config field not settable")
	}
	// Conf value must exist
	confValue, exist := conf[lookup]
	if !exist {
		return errors.Errorf("config value `%s` not found", lookup)
	}
	// Injecting
	objType := obj.Type()
	otk := objType.Kind()
	switch otk {
	case reflect.Int:
		v, err := convertToInt(confValue)
		if err != nil {
			return errors.Wrap(err, "int conversion failed, invalid config value")
		}
		obj.SetInt(int64(v))
	case reflect.String:
		v, err := convertToString(confValue)
		if err != nil {
			return errors.Wrap(err, "string conversion failed, invalid config value")
		}
		obj.SetString(v)
	case reflect.Bool:
		v, err := convertToBool(confValue)
		if err != nil {
			return errors.Wrap(err, "bool conversion failed, invalid config value")
		}
		obj.SetBool(v)
	default:
		return errors.Errorf("invalid config value's type %v", otk)
	}
	return nil
}

func convertToInt(i interface{}) (int, error) {
	switch v := i.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case int32:
		return int(v), nil
	case int16:
		return int(v), nil
	case int8:
		return int(v), nil
	case string:
		v2, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.Wrap(err, "string to int conversion failed")
		}
		return v2, nil
	default:
		return 0, errors.Errorf("type not supported: +%#v", i)
	}
}

func convertToString(i interface{}) (string, error) {
	switch v := i.(type) {
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.Itoa(int(v)), nil
	case int32:
		return strconv.Itoa(int(v)), nil
	case int16:
		return strconv.Itoa(int(v)), nil
	case int8:
		return strconv.Itoa(int(v)), nil
	case string:
		return v, nil
	default:
		return "", errors.Errorf("type not supported: +%#v", i)
	}
}

func convertToBool(i interface{}) (bool, error) {
	switch v := i.(type) {
	case string:
		if v == "true" {
			return true, nil
		}
		if v == "false" {
			return false, nil
		}
		return false, errors.Errorf("invalid value, bool as string expected: +%#v", i)
	case bool:
		return v, nil
	default:
		return false, errors.Errorf("type not supported: +%#v", i)
	}
}
