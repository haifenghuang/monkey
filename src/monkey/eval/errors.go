package eval

import "fmt"

// constants for error types
const (
	_ int = iota
	PREFIXOP
	INFIXOP
	POSTFIXOP
	MOD_ASSIGNOP
	UNKNOWNIDENT
	NOMETHODERROR
	NOINDEXERROR
	KEYERROR
	INDEXERROR
	SLICEERROR
	ARGUMENTERROR
	INPUTERROR
	RTERROR
	PARAMTYPEERROR
	INLENERR
	INVALIDARG
	DIVIDEBYZERO
	THROWERROR
	THROWNOTHANDLED
	NOTLISTABLE
	NOTITERABLE
	RANGETYPEERROR
	DEFERERROR
	SPAWNERROR
	ASSERTIONERROR
	//	STDLIBERROR
	NULLABLEERROR
	JSONERROR
	DBSCANERROR
	FUNCCALLBACKERROR
	FILEMODEERROR
	FILEOPENERROR
	GENERICERROR
)

var errorType = map[int]string{
	PREFIXOP:        "unsupported operator for prefix expression:'%s' and type: %s",
	INFIXOP:         "unsupported operator for infix expression: %s '%s' %s",
	POSTFIXOP:       "unsupported operator for postfix expression:'%s' and type: %s",
	MOD_ASSIGNOP:    "unsupported operator for modulor assignment:'%s'",
	UNKNOWNIDENT:    "unknown identifier: '%s' is not defined",
	NOMETHODERROR:   "undefined method '%s' for object %s",
	NOINDEXERROR:    "index error: type %s is not indexable",
	KEYERROR:        "key error: type %s is not hashable",
	INDEXERROR:      "index error: '%d' out of range",
	SLICEERROR:      "index error: slice '%d:%d' out of range",
	ARGUMENTERROR:   "wrong number of arguments. expected=%s, got=%d",
	INPUTERROR:      "unsupported input type '%s' for function or method: %s",
	RTERROR:         "return type should be %s",
	PARAMTYPEERROR:  "%s argument for '%s' should be type %s. got=%s",
	INLENERR:        "function %s takes input with max length %s. got=%s",
	INVALIDARG:      "invalid argument supplied",
	DIVIDEBYZERO:    "divide by zero",
	THROWERROR:      "throw object must be a string",
	THROWNOTHANDLED: "throw object '%s' not handled",
	NOTLISTABLE:     "grep/map's operating type must be listable",
	NOTITERABLE:     "foreach's operating type must be iterable",
	RANGETYPEERROR:  "range(..) type should be %s type, got='%s'",
	DEFERERROR:      "defer outside function or defer statement not a function",
	SPAWNERROR:      "spawn must be followed by a function",
	ASSERTIONERROR:  "assertion failed",
	//	STDLIBERROR:     "calling '%s' failed",
	NULLABLEERROR:     "%s is null",
	JSONERROR:         "json error: maybe unsupported type or invalid data",
	DBSCANERROR:       "scan type not supported",
	FUNCCALLBACKERROR: "callback error: must be '%d' parameter(s), got '%d'",
	FILEMODEERROR:     "known file mode supplied",
	FILEOPENERROR:     "file open failed, reason: %s",
	GENERICERROR:      "%s",
}

func NewError(line string, t int, args ...interface{}) Object {
	msg := line + fmt.Sprintf(errorType[t], args...)
	return &Error{Kind: t, Message: msg}
}

type Error struct{
	Kind int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func (e *Error) Inspect() string  { return "Runtime Error:" + e.Message + "\n" }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	//	return NewError(line, NOMETHODERROR, method, e.Type())
	return NewError(line, GENERICERROR, e.Message)
}
