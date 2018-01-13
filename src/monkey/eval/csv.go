package eval

import (
	"encoding/csv"
	"os"
)

const (
	CSV_OBJ = "CSV_OBJ"
	csv_name = "csv"
)

var optionsKeyMap = map[string]bool{
	//used in Reader and Writer
	"Comma": true,

	//used in Reader only
	"Comment": true,
	"FieldsPerRecord": true,
	"LazyQuotes": true,
	"TrimLeadingSpace": true,
	"ReuseRecord": true,

	//used in Writer only
	"UseCRLF": true,
}

type CsvObj struct {
	Reader     *csv.Reader
	ReaderFile *os.File

	Writer     *csv.Writer
}

func (c *CsvObj) Inspect() string  { return "<" + csv_name + ">"}
func (c *CsvObj) Type() ObjectType { return CSV_OBJ }

func (c *CsvObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "read":
		return c.Read(line, args...)
	case "readAll":
		return c.ReadAll(line, args...)
	case "closeReader", "close":
		return c.CloseReader(line, args...)
	case "write":
		return c.Write(line, args...)
	case "writeAll":
		return c.WriteAll(line, args...)
	case "flush":
		return c.Flush(line, args...)
	case "setOptions":
		return c.SetOptions(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, c.Type()))
}

func (c *CsvObj) Read(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	record, err  := c.Reader.Read()
	if err != nil {
		return NewNil(err.Error())
	}

	r := &Array{}
	for _, v := range record {
		r.Members = append(r.Members, NewString(v))
	}

	return r
}

func (c *CsvObj) ReadAll(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	records, err  := c.Reader.ReadAll()
	if err != nil {
		return NewNil(err.Error())
	}

	r := &Array{}
	for _, v1 := range records {
		tmpArr := &Array{}
		for _, v2 := range v1 {
			tmpArr.Members = append(tmpArr.Members, NewString(v2))
		}
		r.Members = append(r.Members, tmpArr)
	}

	return r
}

func (c *CsvObj) CloseReader(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if c.ReaderFile != nil {
		err := c.ReaderFile.Close()
		if err != nil {
			return NewNil(err.Error())
		}
	}
	return NIL
}

func (c *CsvObj) Write(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	record, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "Write", "*Array", args[0].Type()))
	}

	strRecord := make([]string, len(record.Members))
	for idx, v := range record.Members {
		strRecord[idx] = v.(*String).String
	}

	err := c.Writer.Write(strRecord)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (c *CsvObj) WriteAll(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	records, ok := args[0].(*Array)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "WriteAll", "*Array", args[0].Type()))
	}

	//create a two dimensional string array
	strRecords := make([][]string, len(records.Members))
	for idx, v := range records.Members {
		strRecords[idx] = make([]string, len(v.(*Array).Members))
	}

	//assign values
	for rowIdx, row := range records.Members {
		cols := row.(*Array)
		for colIdx, col := range cols.Members {
			strRecords[rowIdx][colIdx] = col.(*String).String
		}
	}

	err := c.Writer.WriteAll(strRecords)
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (c *CsvObj) Flush(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	c.Writer.Flush()
	err := c.Writer.Error()
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}

func (c *CsvObj) SetOptions(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	options, ok := args[0].(*Hash)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "setOptions", "*Hash", args[0].Type()))
	}

	for _, option := range options.Pairs {
		//check key type
		key, ok := option.Key.(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "setOptions", "key", "*String", option.Key.Type()))
		}

		//check key name
		if _, ok := optionsKeyMap[key.String]; !ok {
			panic(NewError(line, GENERICERROR, "Keys should be:	Comma|Comment|FieldsPerRecord|LazyQuotes|TrimLeadingSpace|ReuseRecord"))
		}

		//check value type
		switch key.String {
		case "Comma":
			value, ok := option.Value.(*String)
			if !ok {
				panic(NewError(line, GENERICERROR, "'Comma' key's value should be a String"))
			}
			if c.Reader != nil {
				c.Reader.Comma = rune(value.String[0])
			} else if c.Writer != nil {
				c.Writer.Comma = rune(value.String[0])
			}
		case "Comment":
			value, ok := option.Value.(*String)
			if !ok {
				panic(NewError(line, GENERICERROR, "'Comment' key's value should be a String"))
			}
			if c.Reader != nil {
				c.Reader.Comment = rune(value.String[0])
			}
		case "FieldsPerRecord":
			value, ok := option.Value.(*Integer)
			if !ok {
				panic(NewError(line, GENERICERROR, "'FieldsPerRecord' key's value should be an Integer"))
			}
			if c.Reader != nil {
				c.Reader.FieldsPerRecord = int(value.Int64)
			}
		case "LazyQuotes":
			value, ok := option.Value.(*Boolean)
			if !ok {
				panic(NewError(line, GENERICERROR, "'LazyQuotes' key's value should be a Boolean"))
			}
			if c.Reader != nil {
				c.Reader.LazyQuotes = value.Bool
			}
		case "TrimLeadingSpace":
			value, ok := option.Value.(*Boolean)
			if !ok {
				panic(NewError(line, GENERICERROR, "'TrimLeadingSpace' key's value should be a Boolean"))
			}
			if c.Reader != nil {
				c.Reader.TrimLeadingSpace = value.Bool
			}
		case "ReuseRecord":
			value, ok := option.Value.(*Boolean)
			if !ok {
				panic(NewError(line, GENERICERROR, "'ReuseRecord' key's value should be a Boolean"))
			}
			if c.Reader != nil {
				c.Reader.ReuseRecord = value.Bool
			}
		case "UseCRLF":
			value, ok := option.Value.(*Boolean)
			if !ok {
				panic(NewError(line, GENERICERROR, "'UseCRLF' key's value should be a Boolean"))
			}
			if c.Writer != nil {
				c.Writer.UseCRLF = value.Bool
			}
		}
	}

	return NIL
}
