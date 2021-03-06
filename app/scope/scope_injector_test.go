package scope

import (
	"testing"
)

func TestSimpleInject(t *testing.T) {
	t.Parallel()
	var object struct {
		SomeString string `tagname:"SomeStringKey"`
		SomeInt    int    `tagname:"SomeIntKey"`
	}
	dataScope := NewDataScope(map[string]interface{}{
		"SomeStringKey": "SomeStringValue",
		"SomeIntKey":    int(11),
	})
	injector := NewScopeInjector("tagname", dataScope)
	injector.InjectTo(&object)
	if object.SomeInt != 11 {
		t.Error("MapInjector didn't inject a int(11) to SomeInt field")
	}
	if object.SomeString != "SomeStringValue" {
		t.Error("MapInjector didn't inject 'SomeStringValue' to SomeString field")
	}
}
