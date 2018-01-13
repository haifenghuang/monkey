package eval

import (
	_ "fmt"
	"sync"
)

const (
	SYNCCOND_OBJ      = "SYNCCOND_OBJ"
	SYNCONCE_OBJ      = "SYNCONCE_OBJ"
	SYNCMUTEX_OBJ     = "SYNCMUTEX_OBJ"
	SYNCRWMUTEX_OBJ   = "SYNCRWMUTEX_OBJ"
	SYNCWAITGROUP_OBJ = "SYNCWAITGROUP_OBJ"
)

//Condition Object
type SyncCondObj struct {
	Cond *sync.Cond
}

func (c *SyncCondObj) Inspect() string  { return "<" + SYNCCOND_OBJ + ">" }
func (c *SyncCondObj) Type() ObjectType { return SYNCCOND_OBJ }

func (c *SyncCondObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "broadcast":
		return c.Broadcast(line, args...)
	case "signal":
		return c.Signal(line, args...)
	case "wait":
		return c.Wait(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, c.Type()))
}

func (c *SyncCondObj) Broadcast(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	c.Cond.Broadcast()
	return NIL
}

func (c *SyncCondObj) Signal(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	c.Cond.Signal()
	return NIL
}

func (c *SyncCondObj) Wait(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	c.Cond.Wait()
	return NIL
}

//Once Object
type SyncOnceObj struct {
	Once *sync.Once
}

func (o *SyncOnceObj) Inspect() string  { return "<" + SYNCONCE_OBJ + ">" }
func (o *SyncOnceObj) Type() ObjectType { return SYNCONCE_OBJ }

func (o *SyncOnceObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "do":
		return o.Do(line, scope, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, o.Type()))
}

func (o *SyncOnceObj) Do(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	block, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "do", "*Function", args[0].Type()))
	}

	paramCount := len(block.Literal.Parameters)
	if paramCount != 0 {
		panic(NewError(line, FUNCCALLBACKERROR, 0, paramCount))
	}

	o.Once.Do(func() {
		doOnceFunc(scope, block)
	})

	return NIL
}

func doOnceFunc(scope *Scope, f *Function) {
	s := NewScope(scope)

	Eval(f.Literal.Body, s)
	return
}

//Mutex Object
type SyncMutexObj struct {
	Mutex *sync.Mutex
}

func (m *SyncMutexObj) Inspect() string  { return "<" + SYNCMUTEX_OBJ + ">" }
func (m *SyncMutexObj) Type() ObjectType { return SYNCMUTEX_OBJ }

func (m *SyncMutexObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "lock":
		return m.Lock(line, args...)
	case "unlock":
		return m.Unlock(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, m.Type()))
}

func (m *SyncMutexObj) Lock(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	m.Mutex.Lock()
	return NIL
}

func (m *SyncMutexObj) Unlock(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	m.Mutex.Unlock()
	return NIL
}

//RWMutex Object
type SyncRWMutexObj struct {
	RWMutex *sync.RWMutex
}

func (rwm *SyncRWMutexObj) Inspect() string  { return "<" + SYNCRWMUTEX_OBJ + ">" }
func (rwm *SyncRWMutexObj) Type() ObjectType { return SYNCRWMUTEX_OBJ }

func (rwm *SyncRWMutexObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "lock":
		return rwm.Lock(line, args...)
	case "rLock":
		return rwm.RLock(line, args...)
	case "rUnlock":
		return rwm.RUnlock(line, args...)
	case "unlock":
		return rwm.Unlock(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, rwm.Type()))
}

func (rwm *SyncRWMutexObj) Lock(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	rwm.RWMutex.Lock()
	return NIL
}

func (rwm *SyncRWMutexObj) RLock(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	rwm.RWMutex.RLock()
	return NIL
}

func (rwm *SyncRWMutexObj) RUnlock(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	rwm.RWMutex.RUnlock()
	return NIL
}

func (rwm *SyncRWMutexObj) Unlock(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	rwm.RWMutex.Unlock()
	return NIL
}

//WaitGroup Oject
type SyncWaitGroupObj struct {
	WaitGroup *sync.WaitGroup
}

func (wg *SyncWaitGroupObj) Inspect() string  { return "<" + SYNCWAITGROUP_OBJ + ">" }
func (wg *SyncWaitGroupObj) Type() ObjectType { return SYNCWAITGROUP_OBJ }

func (wg *SyncWaitGroupObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "add":
		return wg.Add(line, args...)
	case "done":
		return wg.Done(line, args...)
	case "wait":
		return wg.Wait(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, wg.Type()))
}

func (wg *SyncWaitGroupObj) Add(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	delta, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "do", "*Integer", args[0].Type()))
	}

	wg.WaitGroup.Add(int(delta.Int64))
	return NIL
}

func (wg *SyncWaitGroupObj) Done(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	wg.WaitGroup.Done()
	return NIL
}

func (wg *SyncWaitGroupObj) Wait(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	wg.WaitGroup.Wait()
	return NIL
}
