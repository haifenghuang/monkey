# File `doc.my`
Table of Contents
=================

* Lets
  * [CENTER](#center)
  * [Parse](#parse)
  * [SIZE](#size)

* Enums
  * [Color](#color)
  * [LogOption](#logoption)

* Functions
  * [Add](#add)
  * [Div](#div)
  * [Mul](#mul)
  * [Sub](#sub)

* Classes
  * [DepartmentValidator](#departmentvalidator)

  * [MinMaxValidator](#minmaxvalidator)

  * [NoSpaceValidator](#nospacevalidator)

  * [Person](#person)

    * Lets
      * [firstName](#firstname)
      * [lastName](#lastname)

    * Functions
      * [Msg](#msg)
      * [init](#init)

  * [Request](#request)

  * [RequestHandler](#requesthandler)

  * [RequestMain](#requestmain)

## Lets

### CENTER
```swift
let CENTER = 15.0
```
This is the constant 'CENTER'.

### Parse
```swift
let Parse = fn(lexer) {
	fmt.Printf("This is the parser part!\n") //just a demo
	fmt.Printf("The parser take a lexer as input, and generate an AST\n")
}
```
This is a demo parser

### SIZE
```swift
let SIZE = 10
```
This is the constant 'SIZE'.

## Enums

### Color
```swift
enum Color{
	GRAY = 5, 
	RED = 0, 
	BLUE = 1, 
	YELLOW = 2, 
	BLACK = 3, 
	WHITE = 4
}
```
Color enum is the color of the image.
This is just for demo purpose.

### LogOption
```swift
enum LogOption{
	LUTC = (1 << 5), 
	LstdFlags = ((1 << 4) | (1 << 5)), 
	Ldate = (1 << 0), 
	Ltime = (1 << 1), 
	Lmicroseconds = (1 << 2), 
	Llongfile = (1 << 3), 
	Lshortfile = (1 << 4)
}
```
LogOption enum is defined for Log options.

## Functions

### Add
```swift
fn Add (x, y) 
```
Add returns the addition of its two parameters

This is a example of calling `Add` with two ints:
```swiftsum = Add(1, 2) //result is 3
```

This is a example of calling `Add` with two strings:
```swiftsum = Add("Hello ", "world") //result is "Hello world"
```

#### Parameters
| Name | Type | Description |
| ---- | ---- | ----------- |
|x|`int`| this is the first parameter |
|y|`int`| this is the second parameter |

#### Returns

- `int` the addition of it's two parameters if both parameters are ints. 

- `string` the string concatenation of it's two parameters if both parameters are strings 

#### Source
```swift
fn Add(x,y) {
	return x + y
}
```

### Div
```swift
fn Div (x, y) 
```
Div returns the division of its two parameters

#### Source
```swift
fn Div(x,y) {
	return x / y
}
```

### Mul
```swift
fn Mul (x, y) 
```
Mul returns the multiply of its two parameters

#### Source
```swift
fn Mul(x,y) {
	return x * y
}
```

### Sub
```swift
fn Sub (x, y) 
```
Sub returns the subtraction of its two parameters

#### Parameters
| Name | Type | Description |
| ---- | ---- | ----------- |
|x|`int`| this is the first parameter |
|y|`int`| this is the second parameter |

#### Returns

- `int` the subtraction of it's two parameters 

#### Source
```swift
fn Sub(x,y) {
	return x - y
}
```

## Classes

### DepartmentValidator
```swift
class @DepartmentValidator{ ... }
```
DepartmentValidator annotation class.
It is used for checking a list of department.

#### Source
```swift
class @DepartmentValidator {
  property Department
}
```

### MinMaxValidator
```swift
class @MinMaxValidator{ ... }
```
MinMaxValidator annotation class
It is used for checking min and max length
The max length has a default value of 10

#### Source
```swift
class @MinMaxValidator {
  property MinLength
  property MaxLength default 10
}
```

### NoSpaceValidator
```swift
class @NoSpaceValidator{ ... }
```
NoSpaceValidator annotation class(it's a marker annotation).
It is used for checking spaces.

#### Source
```swift
class @NoSpaceValidator {}
```

### Person
```swift
class Person{ ... }
```
 A Persion class which has three private variable:
* firstName
* lastName
* nickName

and two methods :'init' and 'Msg'.

#### Lets

##### firstName
```swift
let firstName; 
```
this is the firstname

##### lastName
```swift
let lastName; 
```
this is the last name

#### Functions

##### Msg
```swift
fn Msg (extra1, extra2, args) 
```
Method 'Msg' demonstrate function default values and variable arguements

#### Parameters
| Name | Type | Description |
| ---- | ---- | ----------- |
|extra1|`string`| this is extra1 |
|extra2|`string`| this is extra2, it has default value "extra2" |
|args|`array`|this is args |

##### init
```swift
fn init (firstname, lastname, nickname) 
```
This is the constructor of the Person class.

#### Source
```swift
class Person {
	//this is the firstname
	let firstName;
	//this is the last name
	let lastName;
	let nickName; // this is the nick name

	//This is the constructor of the Person class.
	fn init(firstname, lastname="huang", nickname="mike") {
		this.firstName = firstname
		this.lastName = lastname
		this.nickName = nickname
	}

	/* Method 'Msg' demonstrate function default values and variable arguements
	 * @param {string} extra1  this is extra1
	 * @param {string} extra2  this is extra2, it has default value "extra2"
	 * @param {array}  args this is args
	*/
	fn Msg(extra1, extra2="extra2", args...) {
		printf("firstName=%s, lastName=%s, nickName=%s\n", firstName, lastName, nickName)
		printf("extra1=%s, extra2=%s\n", extra1, extra2)
		for i in args {
			printf("i=%v\n", i)
		}
	}
}
```

### Request
```swift
class Request{ ... }
```
Request class represent a user made request.

#### Source
```swift
class Request {
  /* FirstName property, it's value's length should be between 1(minlenght) and 10(maxlength). */
  @MinMaxValidator(MinLength=1)
  property FirstName;

  //LastName property, it's value should contain no space.
  @NoSpaceValidator
  property LastName;

  //Dept property, it's value should be in a list ["Department of Education", "Department of Labors"].
  @DepartmentValidator(Department=["Department of Education", "Department of Labors"])
  property Dept;
}
```

### RequestHandler
```swift
class RequestHandler{ ... }
```
RequestHandler class is responsible for processing the annotation.

#### Source
```swift
class RequestHandler {
  //handle is the real logic of processing the annotations.
  static fn handle(o) {
    props = o.getProperties()
    for p in props {
      annos = p.getAnnotations()
      for anno in annos {
        if anno.instanceOf(MinMaxValidator) {
          //p.value is the property real value.
          if len(p.value) > anno.MaxLength || len(p.value) < anno.MinLength {
            printf("Property '%s' is not valid!\n", p.name)
          }
        } elseif anno.instanceOf(NoSpaceValidator) {
          for c in p.value {
            if c == " " || c == "\t" {
              printf("Property '%s' is not valid!\n", p.name)
              break
            }
          }
        } elseif anno.instanceOf(DepartmentValidator) {
          found = false
          for d in anno.Department {
            if p.value == d {
              found = true
            }
          }
          if !found {
            printf("Property '%s' is not valid!\n", p.name)
          }
        }
      }
    }
  }
}
```

### RequestMain
```swift
class RequestMain{ ... }
```
RequestMain class. It is the user/client code.
#### Note
This is note one.</br>
This is note two.</br>
This is note three.</br>
</br>

#### Warning
This is warning one.</br>
This is warning two.</br>
This is warning three.</br>
</br>

#### Source
```swift
class RequestMain {
  static fn main() {
    request = new Request()
    request.FirstName = "Haifeng123456789"
    request.LastName = "Huang     "
    request.Dept = "Department of Labors"
    RequestHandler.handle(request)
  }
}
```

***
_Last updated 2018-03-05_
