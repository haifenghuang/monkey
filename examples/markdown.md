# File `doc.my`
Table of Contents
=================

* Enums
  * [LogOption](#logoption)

* Classes
  * [Person](#person)

    * Lets
      * [firstName](#firstname)
      * [lastName](#lastname)

    * Functions
      * [Msg](#msg)
      * [init](#init)

## Enums

### LogOption
```monkey
enum LogOption{
	Llongfile = (1 << 3), 
	Lshortfile = (1 << 4), 
	LUTC = (1 << 5), 
	LstdFlags = ((1 << 4) | (1 << 5)), 
	Ldate = (1 << 0), 
	Ltime = (1 << 1), 
	Lmicroseconds = (1 << 2)
}
```
The LogOption enum is defined for Log options.

## Classes

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

***
_Last updated 2018-02-25_
