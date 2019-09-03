package injector

import (
	"reflect"
	"strings"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// MapInjector is map data injector
type MapInjector struct {
	data    map[string]interface{}
	tagname string
}

// NewMapInjector create new map injector instance
func NewMapInjector(tagname string, data map[string]interface{}) app.Injector {
	return app.Injector(MapInjector{
		tagname: tagname,
		data:    data,
	})
}

// InjectTo inject data from all injectors
func (mi MapInjector) InjectTo(obj interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		var isRequired = true
		valueField := structValue.Field(i)
		structField := structValue.Type().Field(i)
		key := structField.Tag.Get(mi.tagname)
		if key == "" {
			continue
		}
		if strings.HasPrefix(key, "?") {
			isRequired = false
			key = key[1:]
		}
		if !valueField.IsValid() {
			return goaterr.Errorf("MapInjector.InjectTo: %s is not valid", structField.Name)
		}
		if !valueField.CanSet() {
			return goaterr.Errorf("MapInjector.InjectTo: Cannot set %s field value", structField.Name)
		}
		newValue, ok := mi.data[key]
		if !ok {
			if !isRequired {
				continue
			}
			return goaterr.Errorf("value for %s is unknown", key)
		}
		if newValue == nil {
			return goaterr.Errorf("MapInjector.InjectTo: dependency instance can not be nil (%s)", key)
		}
		refValue := reflect.ValueOf(newValue)
		valueField.Set(refValue)
	}
	return nil
}
