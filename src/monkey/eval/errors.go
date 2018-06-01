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
	UNKNOWNIDENTEX
	NOMETHODERROR
	NOMETHODERROREX
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
	GREPMAPNOTITERABLE
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
	NOTCLASSERROR
	PARENTNOTDECL
	CLSNOTDEFINE
	CLSMEMBERPRIVATE
	CLSCALLPRIVATE
	PROPERTYUSEERROR
	MEMBERUSEERROR
	INDEXERUSEERROR
	INDEXERTYPEERROR
	INDEXERSTATICERROR
	INDEXNOTFOUNDERROR
	CALLNONSTATICERROR
	CLASSCATEGORYERROR
	PARENTNOTANNOTATION
	OVERRIDEERROR
	METAOPERATORERROR
	GENERICERROR
)

var errorType = map[int]string{
	PREFIXOP:        "unsupported operator for prefix expression:'%s' and type: %s",
	INFIXOP:         "unsupported operator for infix expression: %s '%s' %s",
	POSTFIXOP:       "unsupported operator for postfix expression:'%s' and type: %s",
	MOD_ASSIGNOP:    "unsupported operator for modulor assignment:'%s'",
	UNKNOWNIDENT:    "unknown identifier: '%s' is not defined",
	UNKNOWNIDENTEX:  "identifier '%s' not found. \n\nDid you mean one of: \n\n  %s\n",
	NOMETHODERROR:   "undefined method '%s' for object %s",
	NOMETHODERROREX: "undefined method '%s' for object '%s'. \n\nDid you mean one of: \n\n  %s\n",
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
	GREPMAPNOTITERABLE:  "grep/map's operating type must be iterable",
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
	NOTCLASSERROR:     "Identifier %s is not a class",
	PARENTNOTDECL:     "Parent class %s not declared",
	CLSNOTDEFINE:      "Class %s not defined",
	CLSMEMBERPRIVATE:  "Variable(%s) of class(%s) is private",
	CLSCALLPRIVATE:    "Method %s() of class(%s) is private",
	PROPERTYUSEERROR:  "Invalid use of Property(%s) of class(%s)",
	MEMBERUSEERROR:    "Invalid use of member(%s) of class(%s)",
	INDEXERUSEERROR:   "Invalid use of Indexer of class(%s)",
	INDEXERTYPEERROR:  "Invalid use of Indexer of class(%s), Only interger type of Indexer is supported",
	INDEXERSTATICERROR:"Invalid use of Indexer of class(%s), Indexer cannot declared as static",
	INDEXNOTFOUNDERROR:"Indexer not found for class(%s)",
	CALLNONSTATICERROR:"Could not call non-static.",
	CLASSCATEGORYERROR:"No class(%s) found for category(%s).",
	PARENTNOTANNOTATION:"Annotation(%s)'s Parent(%s) is not annotation.",
	OVERRIDEERROR:      "Method(%s) of class(%s) must override a superclass method!",
	METAOPERATORERROR:  "Meta-Operators' item must be Numbers|String!",
	GENERICERROR:      "%s",
}

func NewError(line string, t int, args ...interface{}) Object {
	msg := line + fmt.Sprintf(errorType[t], args...)
	return &Error{Kind: t, Message: msg}
}

type Error struct {
	Kind    int
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
