package builder

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/omcrgnt/res"
)

// SeedMap maps first-level AppResources field name to config spec in reg.
type SeedMap map[string]any

var seedMaps sync.Map

// newResourceSpec defers [NewResourceer.NewResource] until [Build].
type newResourceSpec struct {
	factory NewResourceer
}

func (s newResourceSpec) Build() (any, error) {
	return s.factory.NewResource()
}

// Seed walks first-level fields of appResources and registers build specs only.
// [BuildConfiger] → spec from [BuildConfiger.BuildConfig]; [NewResourceer] → [newResourceSpec].
func Seed(reg res.Registry, appResources any) error {
	if reg == nil {
		return fmt.Errorf("builder: nil registry")
	}

	rv, err := structValue(appResources)
	if err != nil {
		return err
	}

	rt := rv.Type()
	m := make(SeedMap)
	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fieldVal := rv.Field(i)
		if !fieldVal.CanInterface() {
			continue
		}

		field := fieldVal.Interface()
		_, isNew := field.(NewResourceer)
		_, isBuild := field.(BuildConfiger)
		switch {
		case isNew && isBuild:
			return fmt.Errorf("builder: %s: implements both NewResourceer and BuildConfiger", sf.Name)
		case isNew:
			spec := newResourceSpec{factory: field.(NewResourceer)}
			if err := reg.Add(spec); err != nil {
				return fmt.Errorf("builder: %s: %w", sf.Name, err)
			}
		case isBuild:
			spec, err := field.(BuildConfiger).BuildConfig()
			if err != nil {
				return fmt.Errorf("builder: %s: BuildConfig: %w", sf.Name, err)
			}
			if err := reg.Add(spec); err != nil {
				return fmt.Errorf("builder: %s: %w", sf.Name, err)
			}
			m[sf.Name] = spec
			if tag, ok := sf.Tag.Lookup("ecfg"); ok {
				if seg := strings.ToUpper(strings.Split(tag, ",")[0]); seg != "" {
					m[seg] = spec
				}
			}
		default:
			return fmt.Errorf("builder: %s: must implement NewResourceer or BuildConfiger", sf.Name)
		}
	}

	seedMaps.Store(reg, m)
	return nil
}

// SeedMapFor returns the map populated by [Seed] for reg.
func SeedMapFor(reg res.Registry) (SeedMap, bool) {
	v, ok := seedMaps.Load(reg)
	if !ok {
		return nil, false
	}
	return v.(SeedMap), true
}

func structValue(v any) (reflect.Value, error) {
	if v == nil {
		return reflect.Value{}, fmt.Errorf("builder: nil app resources")
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, fmt.Errorf("builder: nil app resources")
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("builder: want struct, got %s", rv.Kind())
	}

	return rv, nil
}
