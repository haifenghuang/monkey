package eval

import (
	"database/sql"
	_ "fmt"
	//	_ "github.com/mattn/go-sqlite3"
	_ "reflect"
)

const (
	DBSQL_OBJ    = "DBSQL_OBJ"    //sql object
	DBRESULT_OBJ = "DBRESULT_OBJ" //result object

	DBROW_OBJ  = "DBROW_OBJ"  //row object
	DBROWS_OBJ = "DBROWS_OBJ" //rows object

	DBSTMT_OBJ = "DBSTMT_OBJ" //statement object
	DBTX_OBJ   = "DBTX_OBJ"   //transaction object
)

const (
	SQL_OBJ = "SQL_OBJ"
	sql_name = "sql"
)

//This object's purpose is only for 5 predefined null constants
type SqlsObject struct {
}

func (s *SqlsObject) Inspect() string  { return "<" + sql_name + ">" }
func (s *SqlsObject) Type() ObjectType { return SQL_OBJ }
func (s *SqlsObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	panic(NewError(line, NOMETHODERROR, method, s.Type()))
}

func NewSqlsObject() Object {
	ret := &SqlsObject{}
	SetGlobalObj(sql_name, ret)

	//Below five variables are mainly used in `sql` handling.
	//when we want to insert null to a certain column, we could
	//use these variables
	SetGlobalObj(sql_name+".INT_NULL",     &Integer{Int64: 0, Valid: false})    //NullInteger
	SetGlobalObj(sql_name+".UINT_NULL",    &UInteger{UInt64: 0, Valid: false})    //NullUInteger
	SetGlobalObj(sql_name+".FLOAT_NULL",   &Float{Float64: 0.0, Valid: false})  //NullFloat
	SetGlobalObj(sql_name+".STRING_NULL",  &String{String: "", Valid: false})   //NullString
	SetGlobalObj(sql_name+".BOOL_NULL",    &Boolean{Bool: true, Valid: false})  //NullBool
	SetGlobalObj(sql_name+".TIME_NULL",    &TimeObj{Valid: false})              //NullTime
	SetGlobalObj(sql_name+".DECIMAL_NULL", &DecimalObj{Valid:false})            //NullDecimal

	return ret
}

//***************************************************************
//                         SQL Object
//***************************************************************
type SqlObject struct {
	Db   *sql.DB
	Name string
}

// Implement the 'Closeable' interface
func (s *SqlObject) close(line string, args ...Object) Object {
	return s.Close(line, args...)
}

func (s *SqlObject) Inspect() string  { return s.Name }
func (s *SqlObject) Type() ObjectType { return DBSQL_OBJ }
func (s *SqlObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "ping":
		return s.Ping(line, args...)
	case "close":
		return s.Close(line, args...)
	case "setMaxOpenConns":
		return s.SetMaxOpenConns(line, args...)
	case "setMaxIdleConns":
		return s.SetMaxIdleConns(line, args...)
	case "exec":
		return s.Exec(line, args...)
	case "query":
		return s.Query(line, args...)
	case "queryRow":
		return s.QueryRow(line, args...)
	case "prepare":
		return s.Prepare(line, args...)
	case "begin":
		return s.Begin(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, s.Type()))
	}
}

//Return the remote address
func (s *SqlObject) Ping(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := s.Db.Ping()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (s *SqlObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := s.Db.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (s *SqlObject) SetMaxOpenConns(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	n, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setMaxOpenConns", "*Integer", args[0].Type()))
	}

	s.Db.SetMaxOpenConns(int(n.Int64))
	return NIL
}

func (s *SqlObject) SetMaxIdleConns(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	n, ok := args[0].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setMaxIdleConns", "*Integer", args[0].Type()))
	}

	s.Db.SetMaxIdleConns(int(n.Int64))
	return NIL
}

func (s *SqlObject) Exec(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "exec", "*String", args[0].Type()))
	}

	var params []interface{}
	if len(args) > 1 {
		params = handleExecParams(args[1:])
		if params == nil {
			panic(NewError(line, INVALIDARG))
		}
	}

	result, err := s.Db.Exec(query.String, params...)
	if err != nil {
		return NewNil(err.Error())
	}

	return &DbResultObject{Result: result, Name: s.Name}
}

func (s *SqlObject) Query(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "query", "*String", args[0].Type()))
	}

	var params []interface{}
	arr := args[1:]
	for idx, arg := range arr {
		_, ok := arg.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "remaining", "query", "*String", arr[idx].Type()))
		}
		params = append(params, arg.Inspect())
	}

	rows, err := s.Db.Query(query.String, params...)
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbRowsObject{Rows: rows, Name: s.Name}

}

func (s *SqlObject) QueryRow(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "queryRow", "*String", args[0].Type()))
	}

	var params []interface{}
	arr := args[1:]
	for idx, arg := range arr {
		_, ok := arg.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "remaining", "queryRow", "*String", arr[idx].Type()))
		}
		params = append(params, arg.Inspect())
	}

	row := s.Db.QueryRow(query.String, params...)
	return &DbRowObject{Row: row, Name: s.Name}
}

func (s *SqlObject) Prepare(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "prepare", "*String", args[0].Type()))
	}

	stmt, err := s.Db.Prepare(query.String)
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbStmtObject{Stmt: stmt, Name: s.Name}
}

func (s *SqlObject) Begin(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbTxObject{Tx: tx, Name: s.Name}
}

//***************************************************************
//                         DB Result Object
//***************************************************************
type DbResultObject struct {
	Result sql.Result
	Name   string
}

func (r *DbResultObject) Inspect() string  { return r.Name }
func (r *DbResultObject) Type() ObjectType { return DBRESULT_OBJ }
func (r *DbResultObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "lastInsertId":
		return r.LastInsertId(line, args...)
	case "rowsAffected":
		return r.RowsAffected(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, r.Type()))
	}
}

func (r *DbResultObject) LastInsertId(line string, args ...Object) Object {
	n, err := r.Result.LastInsertId()
	if err != nil {
		return NewInteger(-1)
	}
	return NewInteger(n)
}

func (r *DbResultObject) RowsAffected(line string, args ...Object) Object {
	n, err := r.Result.RowsAffected()
	if err != nil {
		return NewInteger(-1)
	}
	return NewInteger(n)
}

//***************************************************************
//                         DB Rows Object
//***************************************************************
type DbRowsObject struct {
	Rows *sql.Rows
	Name string
}

// Implement the 'Closeable' interface
func (r *DbRowsObject) close(line string, args ...Object) Object {
	return r.Close(line, args...)
}

func (r *DbRowsObject) Inspect() string  { return r.Name }
func (r *DbRowsObject) Type() ObjectType { return DBROWS_OBJ }
func (r *DbRowsObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "columns":
		return r.Columns(line, args...)
	case "scan":
		return r.Scan(line, args...)
	case "next":
		return r.Next(line, args...)
	case "close":
		return r.Close(line, args...)
	case "err":
		return r.Err(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, r.Type()))
	}
}

func (r *DbRowsObject) Columns(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	arr := &Array{}
	cols, err := r.Rows.Columns()
	if err != nil {
		return NewNil(err.Error())
	}

	for _, col := range cols {
		arr.Members = append(arr.Members, NewString(col))
	}

	return arr
}

func (r *DbRowsObject) Scan(line string, args ...Object) Object {
	return scan(r.Rows, line, args...)
}

func (r *DbRowsObject) Next(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	b := r.Rows.Next()
	if b {
		return TRUE
	}
	return FALSE
}

func (r *DbRowsObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := r.Rows.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (r *DbRowsObject) Err(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := r.Rows.Err()
	if err != nil {
		return NewString(err.Error())
	}

	return NewString("")
}

//***************************************************************
//                         DB Row Object
//***************************************************************
type DbRowObject struct {
	Row  *sql.Row
	Name string
}

func (r *DbRowObject) Inspect() string  { return r.Name }
func (r *DbRowObject) Type() ObjectType { return DBROW_OBJ }
func (r *DbRowObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "scan":
		return r.Scan(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, r.Type()))
	}
}

func (r *DbRowObject) Scan(line string, args ...Object) Object {
	return scan(r.Row, line, args...)
}

//***************************************************************
//                         DB Statement Object
//***************************************************************
type DbStmtObject struct {
	Stmt *sql.Stmt
	Name string
}

// Implement the 'Closeable' interface
func (s *DbStmtObject) close(line string, args ...Object) Object {
	return s.Close(line, args...)
}

func (s *DbStmtObject) Inspect() string  { return s.Name }
func (s *DbStmtObject) Type() ObjectType { return DBSTMT_OBJ }
func (s *DbStmtObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "exec":
		return s.Exec(line, args...)
	case "query":
		return s.Query(line, args...)
	case "queryRow":
		return s.QueryRow(line, args...)
	case "close":
		return s.Close(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, s.Type()))
	}
}
func (s *DbStmtObject) Close(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}
	err := s.Stmt.Close()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (s *DbStmtObject) Exec(line string, args ...Object) Object {

	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	params := handleExecParams(args)
	if params == nil {
		panic(NewError(line, INVALIDARG))
	}

	result, err := s.Stmt.Exec(params...)
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbResultObject{Result: result, Name: s.Name}
}

func (s *DbStmtObject) Query(line string, args ...Object) Object {
	var params []interface{}
	for idx, arg := range args {
		_, ok := arg.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "All", "query", "*String", args[idx].Type()))
		}
		params = append(params, arg.Inspect())
	}

	rows, err := s.Stmt.Query(params...)
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbRowsObject{Rows: rows, Name: s.Name}

}

func (s *DbStmtObject) QueryRow(line string, args ...Object) Object {
	var params []interface{}
	for idx, arg := range args {
		_, ok := arg.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "all", "queryRow", "*String", args[idx].Type()))
		}
		params = append(params, arg.Inspect())
	}

	row := s.Stmt.QueryRow(params...)
	return &DbRowObject{Row: row, Name: s.Name}
}

//***************************************************************
//                         DB Transaction object
//***************************************************************
type DbTxObject struct {
	Tx   *sql.Tx
	Name string
}

func (t *DbTxObject) Inspect() string  { return t.Name }
func (t *DbTxObject) Type() ObjectType { return DBTX_OBJ }
func (t *DbTxObject) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "exec":
		return t.Exec(line, args...)
	case "query":
		return t.Query(line, args...)
	case "queryRow":
		return t.QueryRow(line, args...)
	case "prepare":
		return t.Prepare(line, args...)
	case "stmt":
		return t.Stmt(line, args...)
	case "commit":
		return t.Commit(line, args...)
	case "rollback":
		return t.Rollback(line, args...)
	default:
		panic(NewError(line, NOMETHODERROR, method, t.Type()))
	}
}

func (t *DbTxObject) Exec(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "exec", "*String", args[0].Type()))
	}

	var params []interface{}
	if len(args) > 1 {
		params = handleExecParams(args[1:])
		if params == nil {
			panic(NewError(line, INVALIDARG))
		}
	}

	result, err := t.Tx.Exec(query.String, params...)
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbResultObject{Result: result, Name: t.Name}
}

func (t *DbTxObject) Query(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "query", "*String", args[0].Type()))
	}

	var params []interface{}
	arr := args[1:]
	for idx, arg := range arr {
		_, ok := arg.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "remaining", "query", "*String", arr[idx].Type()))
		}
		params = append(params, arg.Inspect())
	}

	rows, err := t.Tx.Query(query.String, params...)
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbRowsObject{Rows: rows, Name: t.Name}

}

func (t *DbTxObject) QueryRow(line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, "at least 1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "queryRow", "*String", args[0].Type()))
	}

	var params []interface{}
	arr := args[1:]
	for idx, arg := range arr {
		_, ok := arg.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "remaining", "queryRow", "*String", arr[idx].Type()))
		}
		params = append(params, arg.Inspect())
	}

	row := t.Tx.QueryRow(query.String, params...)
	return &DbRowObject{Row: row, Name: t.Name}
}

func (t *DbTxObject) Prepare(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	query, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "prepare", "*String", args[0].Type()))
	}

	stmt, err := t.Tx.Prepare(query.String)
	if err != nil {
		return NewNil(err.Error())
	}
	return &DbStmtObject{Stmt: stmt, Name: t.Name}
}

func (t *DbTxObject) Stmt(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	stmt, ok := args[0].(*DbStmtObject)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "stmt", "*DbStmtObject", args[0].Type()))
	}

	newStmt := t.Tx.Stmt(stmt.Stmt)
	return &DbStmtObject{Stmt: newStmt, Name: t.Name}
}

func (t *DbTxObject) Commit(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := t.Tx.Commit()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

func (t *DbTxObject) Rollback(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	err := t.Tx.Rollback()
	if err != nil {
		return NewFalseObj(err.Error())
	}
	return TRUE
}

//Handling `Exec`'s parameters, mainly for handling `null`
//func handleExecParams(args []Object) []interface{} {
//
//	var params []interface{}
//	for _, arg := range args {
//		switch arg.(type) {
//		case *String:
//			var nulStr sql.NullString
//
//			s := arg.(*String)
//			if s.Valid {
//				nulStr = sql.NullString{String:s.String, Valid:true}
//			} else {
//				nulStr = sql.NullString{}
//			}
//			params = append(params, nulStr)
//		case *Integer:
//			var nulInt sql.NullInt64
//
//			i := arg.(*Integer)
//			if i.Valid {
//				nulInt = sql.NullInt64{Int64:i.Int64, Valid:true}
//			} else {
//				nulInt = sql.NullInt64{}
//			}
//			params = append(params, nulInt)
//		case *Float:
//			var nulFloat sql.NullFloat64
//
//			f := arg.(*Float)
//			if f.Valid {
//				nulFloat = sql.NullFloat64{Float64:f.Float64, Valid:true}
//			} else {
//				nulFloat = sql.NullFloat64{}
//			}
//			params = append(params, nulFloat)
//		case *Boolean:
//			var nulBool sql.NullBool
//
//			b := arg.(*Boolean)
//			if b.Valid {
//				nulBool = sql.NullBool{Bool:b.Bool, Valid:true}
//			} else {
//				nulBool = sql.NullBool{}
//			}
//			params = append(params, nulBool)
//		default:
//			return nil
//		} //end switch
//	}
//
//	return params
//}

//Handling `Exec`'s parameters, mainly for handling `null`
func handleExecParams(args []Object) []interface{} {

	var params []interface{}
	for _, arg := range args {
		params = append(params, arg)
	}

	return params
}

func scan(rowObj interface{}, line string, args ...Object) Object {
	if len(args) < 1 {
		panic(NewError(line, ARGUMENTERROR, ">=1", len(args)))
	}

	//Must do `[]Object ==> []interface{}`，or compiler will report error
	var values []interface{}
	for _, v := range args {
		switch v.(type) {
		case *Integer, *UInteger, *Boolean, *Float, *String, *TimeObj:
			values = append(values, v)
		default:
			panic(NewError(line, DBSCANERROR))
		} //end switch
	} //end for

	var err error
	switch rowObj.(type) {
	case *sql.Rows:
		row := rowObj.(*sql.Rows)
		err = row.Scan(values...)
	case *sql.Row:
		row := rowObj.(*sql.Row)
		err = row.Scan(values...)
	}

	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}
