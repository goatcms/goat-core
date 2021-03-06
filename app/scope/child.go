package scope

import (
	"context"
	"sync"

	"github.com/goatcms/goatcore/app"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// ChildScope represent sub-scope
type ChildScope struct {
	app.DataScope
	app.EventScope
	parent    app.Scope
	errors    []error
	errorsMU  sync.Mutex
	waitGroup sync.WaitGroup
	injectors []app.Injector
}

// NewChildScope create new instance of scope
func NewChildScope(parent app.Scope, params ChildParams) app.Scope {
	parent.AddTasks(1)
	if params.DataScope == nil {
		params.DataScope = NewChildDataScope(parent, make(map[string]interface{}))
	}
	if params.EventScope == nil {
		params.EventScope = NewChildEventScope(parent)
	}
	return &ChildScope{
		parent:     parent,
		DataScope:  params.DataScope,
		EventScope: params.EventScope,
		injectors:  params.Injectors,
	}
}

// Context return shared context
func (cs *ChildScope) Context() context.Context {
	return cs.parent.Context()
}

// Kill shared context
func (cs *ChildScope) Kill() {
	cs.parent.Kill()
}

// IsKilled return true if shared context is killed/ended
func (cs *ChildScope) IsKilled() bool {
	return cs.parent.IsKilled()
}

// Wait for end of all tasks in child scope
func (cs *ChildScope) Wait() (err error) {
	cs.waitGroup.Wait()
	return goaterr.ToError(cs.errors)
}

// AddTasks tasks to child scope
func (cs *ChildScope) AddTasks(delta int) (err error) {
	cs.waitGroup.Add(delta)
	return nil
}

// DoneTask done one child scope task
func (cs *ChildScope) DoneTask() {
	cs.waitGroup.Done()
}

// Errors return child scope errors
func (cs *ChildScope) Errors() []error {
	return cs.errors
}

// ToError return error if child scope contains any error
func (cs *ChildScope) ToError() error {
	return goaterr.ToError(cs.errors)
}

// AppendError add error to child and parent scope
func (cs *ChildScope) AppendError(err error) {
	cs.errorsMU.Lock()
	defer cs.errorsMU.Unlock()
	cs.Kill()
	cs.parent.AppendError(err)
	cs.errors = append(cs.errors, err)
}

// AppendErrors add errors to child and parent scope
func (cs *ChildScope) AppendErrors(errs ...error) {
	cs.errorsMU.Lock()
	defer cs.errorsMU.Unlock()
	cs.Kill()
	cs.parent.AppendErrors(errs...)
	cs.errors = append(cs.errors, errs...)
}

// InjectTo insert data to object
func (cs *ChildScope) InjectTo(obj interface{}) (err error) {
	for _, scpInjector := range cs.injectors {
		if err = scpInjector.InjectTo(obj); err != nil {
			return err
		}
	}
	return cs.parent.InjectTo(obj)
}

// Close child scope
func (cs *ChildScope) Close() (err error) {
	err = cs.Wait()
	cs.parent.DoneTask()
	return err
}
