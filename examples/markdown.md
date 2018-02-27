# File `doc.my`
Table of Contents
=================

* Lets
  * [CENTER](#center)
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
let CENTER = 15.0; 
```
This is the constant 'CENTER'.

### SIZE
```swift
let SIZE = 10; 
```
This is the constant 'SIZE'.

## Enums

### Color
```swift
enum Color{
	WHITE = 4, 
	GRAY = 5, 
	RED = 0, 
	BLUE = 1, 
	YELLOW = 2, 
	BLACK = 3
}
```
Color enum is the color of the image.
This is just for demo purpose.

### LogOption
```swift
enum LogOption{
	Lmicroseconds = (1 << 2), 
	Llongfile = (1 << 3), 
	Lshortfile = (1 << 4), 
	LUTC = (1 << 5), 
	LstdFlags = ((1 << 4) | (1 << 5)), 
	Ldate = (1 << 0), 
	Ltime = (1 << 1)
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
```swift
sum = Add(1, 2) //result is 3
```

This is a example of calling `Add` with two strings:
```swift
sum = Add("Hello ", "world") //result is "Hello world"
```

#### Parameters
| Name | Type | Description |
| ---- | ---- | ----------- |
x|`int`| this is the first parameter |
y|`int`| this is the second parameter |

#### Returns

- `int` the addition of it's two parameters if both parameters are ints. 

- `string` the string concatenation of it's two parameters if both parameters are strings 

### Div
```swift
fn Div (x, y) 
```
Div returns the division of its two parameters

### Mul
```swift
fn Mul (x, y) 
```
Mul returns the multiply of its two parameters

### Sub
```swift
fn Sub (x, y) 
```
Sub returns the subtraction of its two parameters

#### Parameters
| Name | Type | Description |
| ---- | ---- | ----------- |
x|`int`| this is the first parameter |
y|`int`| this is the second parameter |

#### Returns

- `int` the subtraction of it's two parameters 

## Classes

### DepartmentValidator
```swift
class @DepartmentValidator{ ... }
```
DepartmentValidator annotation class.
It is used for checking a list of department.

### MinMaxValidator
```swift
class @MinMaxValidator{ ... }
```
MinMaxValidator annotation class
It is used for checking min and max length
The max length has a default value of 10

### NoSpaceValidator
```swift
class @NoSpaceValidator{ ... }
```
NoSpaceValidator annotation class(it's a marker annotation).
It is used for checking spaces.

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
extra1|`string`| this is extra1 |
extra2|`string`| this is extra2, it has default value "extra2" |
args|`array`|this is args |

##### init
```swift
fn init (firstname, lastname, nickname) 
```
This is the constructor of the Person class.

### Request
```swift
class Request{ ... }
```
Request class represent a user made request.

### RequestHandler
```swift
class RequestHandler{ ... }
```
RequestHandler class is responsible for processing the annotation.

### RequestMain
```swift
class RequestMain{ ... }
```
RequestMain class. It is the user/client code.

***
_Last updated 2018-02-27_
