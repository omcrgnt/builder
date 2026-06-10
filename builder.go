package builder

import (
	"fmt"
	"reflect"
)

// Builder — поле cfg (первого уровня), способное собрать ресурс приложения.
type Builder interface {
	Build() (any, error)
}

// Registrar принимает готовые ресурсы (res реализует через Add / AddAll).
type Registrar interface {
	Add(any) error
}

// Build обходит поля v первого уровня, вызывает Build() у реализующих Builder
// и регистрирует результат. При первой ошибке возвращает её с именем поля.
func Build(v any, reg Registrar) error {
	if reg == nil {
		return fmt.Errorf("builder: nil registrar")
	}

	rv, err := structValue(v)
	if err != nil {
		return err
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fieldVal := rv.Field(i)
		if !fieldVal.CanInterface() {
			continue
		}
		if isNilValue(fieldVal) {
			continue
		}

		b, ok := fieldVal.Interface().(Builder)
		if !ok {
			continue
		}

		res, err := b.Build()
		if err != nil {
			return fmt.Errorf("builder: %s: %w", rt.Field(i).Name, err)
		}

		if err := reg.Add(res); err != nil {
			return fmt.Errorf("builder: %s: %w", rt.Field(i).Name, err)
		}
	}

	return nil
}

func structValue(v any) (reflect.Value, error) {
	if v == nil {
		return reflect.Value{}, fmt.Errorf("builder: nil config")
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, fmt.Errorf("builder: nil config")
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("builder: want struct, got %s", rv.Kind())
	}

	return rv, nil
}

func isNilValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}
