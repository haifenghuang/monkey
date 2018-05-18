package eval

import (
	"fmt"
)

type ChanObject struct {
	ch chan Object
	done bool
}

//Make channel object could be used in `for x in channelObj`
func (c *ChanObject) iter() bool { return true }

//Implement the 'Closeable' interface
func (c *ChanObject) close(line string, args ...Object) Object {
	return c.Close(line, args...)
}

func (c *ChanObject) Inspect() string  { return fmt.Sprintf("channel<%p>", c.ch) }
func (c *ChanObject) Type() ObjectType { return CHANNEL_OBJ }
func (c *ChanObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "send":
		return c.Send(line, args...)
	case "recv":
		return c.Recv(line, args...)
	case "close":
		return c.Close(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, c.Type()))
	}
}

func (c *ChanObject) Send(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	c.ch <- args[0]
	return NIL
}

func (c *ChanObject) Recv(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	obj, more := <-c.ch
	c.done = more
	return obj
}

func (c *ChanObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	close(c.ch)
	return NIL
}
