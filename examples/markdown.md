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
```monkey
let CENTER = 15.0; 
```
This is the constant 'CENTER'.

### SIZE
```monkey
let SIZE = 10; 
```
This is the constant 'SIZE'.

## Enums

### Color
```monkey
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
```monkey
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
```monkey
fn Add (x, y) 
```
Add returns the sum of its two parameters

### Div
```monkey
fn Div (x, y) 
```
Div returns the division of its two parameters

### Mul
```monkey
fn Mul (x, y) 
```
Mul returns the multiply of its two parameters

### Sub
```monkey
fn Sub (x, y) 
```
Sub returns the subtraction of its two parameters

## Classes

### DepartmentValidator
```monkey
class @DepartmentValidator{ ... }
```
DepartmentValidator annotation class.
 It is used for checking a list of department.

### MinMaxValidator
```monkey
class @MinMaxValidator{ ... }
```
MinMaxValidator annotation class
It is used for checking min and max length
The max length has a default value of 10

### NoSpaceValidator
```monkey
class @NoSpaceValidator{ ... }
```
NoSpaceValidator annotation class(it's a marker annotation).
It is used for checking spaces.

### Person
```monkey
class Person{ ... }
```
  A Persion class which has three private variable:
       firstName
       lastName
       nickName
  and two methods :'init' and 'Msg'.

#### Lets

##### firstName
```monkey
let firstName; 
```
this is the firstname

##### lastName
```monkey
let lastName; 
```
this is the last name

#### Functions

##### Msg
```monkey
fn Msg (extra1, extra2, args) 
```
Method 'Msg' demonstrate function default values and variable arguements

##### init
```monkey
fn init (firstname, lastname, nickname) 
```
This is the constructor of the Person class.

### Request
```monkey
class Request{ ... }
```
Request class represent a user made request.

### RequestHandler
```monkey
class RequestHandler{ ... }
```
RequestHandler class is responsible for processing the annotation.

### RequestMain
```monkey
class RequestMain{ ... }
```
RequestMain class. It is the user/client code.

***
_Last updated 2018-02-26_
