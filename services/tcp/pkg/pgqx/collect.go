package pgqx

import (
	"errors"
	"reflect"
)

func collectFields(target any) (map[string]any, error) {
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return nil, errors.New("target must be non-nil pointer")
	}

	value = value.Elem()
	t := value.Type()
	fields := make(map[string]any)

	for i := 0; i < value.NumField(); i++ {
		fv := value.Field(i)
		ft := t.Field(i)

		if ft.PkgPath != "" {
			continue
		}

		dbTag := ft.Tag.Get("db")
		if dbTag == "" {
			dbTag = ft.Name
		}

		if fv.Kind() == reflect.Struct && ft.Type.String() != "time.Time" {
			nested, err := collectFields(fv.Addr().Interface())
			if err != nil {
				return nil, err
			}

			for k, v := range nested {
				fields[k] = v
			}

			continue
		}

		fields[dbTag] = fv.Addr().Interface()
	}

	return fields, nil
}
