# Monkey Programming Language

Chinese version: [中文](README_cn.md)

## Summary

Monkey is a toy language interpreter, written in Go. It has C-style syntax, and is largely inspired by Ruby, Python and Perl.

## Overview

This project is based on mayoms's project [monkey](https://github.com/mayoms/monkey) with some bug fixes and a lot of new features including:

* Modified string module(which can correctly handle utf8 character encoding)
* Added file module(with some new methods).
* Added math module
* Added sql(db) module(which can correctly handing null values. Need sql driver package)
* Added time module
* Added sort module
* Added os module
* Added log module
* Added net module
* Added http module
* Added filepath module
* Added flag module(for handling command line options)
* Added json module(for json marshaling and unmarshaling)
* Added fmt module
* Added sync module
* Added list module
* Added linq module(Code come from [linq](https://github.com/ahmetb/go-linq) with some modifications)
* Added csv module
* Added regexp module
* Added template module
* Regular expression literal support(partially like perls)
* channel support(like golang's channel)
* more operator support(&&, ||, &, |, ^, +=, -=, ?: etc.)
* utf8 support(e.g. you could use utf8 character as variable name)
* more flow control support(e.g. try/catch/finally, for-in, case, c-like for loop)
* defer support
* spawn support(goroutine)
* enum support
* pipe operator support(see demo for help)
* function with default value and variadic parameters

There are a number of tasks to complete, as well as a number of bugs. The purpose of this project was to dive deeper into Go, as well as get a better understanding of how programming languages work. It has been successful in those goals. There may or may not be continued work - I do plan on untangling a few messy spots, and there are a few features I'd like to see implemented. This will happen as time and interest allows.

## Installation

Just download the repository and run `./run.sh`

## Basic use

To access the REPL, simply run the following:

```sh
~ » monkey
Monkey programming language REPL

>>
```

or, to run a program:

```sh
monkey path/to/file
```

## Language Tour

### Comments

Monkey only support single line comment.

```swift
// an inline comment
# another inline comment
```

### Data Types

Monkey supports 7 basic data types: `String`, `Int`, `Float`, `Bool`, `Array`, `Hash` and `Nil`

```swift
s1 = "hello, 黄"      # strings are UTF-8 encoded
s2 = `hello, "world"`  # raw string
i = 10                 # int
f = 10.0               # float
b = true               # bool
a = [1, "2"]           # array
h = { "a"=>1, "b"=>2}  # hash
n = nil
```

### Constants(Literal)

In monkey, there are mainly nine types of constants(Literals).

* Integer
* Float
* String
* Regular expression
* Array
* Hash
* Nil
* Boolean
* Function

```swift
// Integer literals
i1 = 10
i2 = 20_000_000     //for more readable
i3 = 0x80           // hex
i4 = 0b10101        // binary
i5 = 0c127          // octal

// Float literals
f1 = 10.25
f2 = 1.02E3
f3 = 123_456.789_012 //for more readable

// String literals
s1 = "123"
s2 = "Hello world"

// Regular expression literals
r = /\d+/.match("12")
if (r) { prinln("regex matched!") }

// Array literals
a = [1+2, 3, 4, "5", 3]

// Hash literals
h = { "a"=>1, "b"=>2, "c"=>2}

// Nil literal
n = nil

// Boolean literals
t = true
f = false

// Function literals
let f = add(x, y) { return a + b }
println(f(1,2))
```

### Variables

Variables in Monkey could start with the keyword `let`, or nothing with the
form `variable=value`.

```swift
let a = 1, b = "hello world", c = [1,2,3]
d = 4
e = 5
姓 = "黄"
```

### Reserved keywords

Keywords are predefined, reserved identifiers that have special meanings to the compiler. They cannot be used as identifiers. Below is a list of reserved keywords

* fn
* let
* true false nil
* if elsif elseif else
* return
* include
* and or
* enum
* struct # not used
* do while for break continue where
* grep map
* case is in
* try catch finally throw
* defer
* spawn
* yield #not used
* qw

### Type conversion

You can use the builtin `int()`, `float()`, `str()`, `array()` functions for type conversion.

```swift
let i = 0xa
let s = str(i)     // result: "10"
let f = float(i)   // result: 10
let a = array(i)   // result: [10]
```

### `qw`(Quote word) keyword

The `qw` keyword is like perl's `qw` keyword. When you want to use a lot of quoted strings, the `qw` keyword can make it a lot easier for those strings.

```swift
for str in qw<abc, def, ghi, jkl, mno> { //allowed 'qw' pair is '{}', '<>', '()'
  println('str={str}')
}

newArr = qw(1,2,3.5) //array with string values, not number values.
fmt.printf("newArr=%v\n", newArr)
```

### `enum` keyword

In mokey, you can use enum to define constants.

```swift
LogOption = enum {
    Ldate         = 1 << 0,
    Ltime         = 1 << 1,
    Lmicroseconds = 1 << 2,
    Llongfile     = 1 << 3,
    Lshortfile    = 1 << 4,
    LUTC          = 1 << 5,
    LstdFlags     = 1 << 4 | 1 << 5
}

opt = LogOption.LstdFlags
println(opt)

//get all names of the `enum`
for s in LogOption.getNames() { //not ordered
    println(s)
}

//get all values of the `enum`
for s in LogOption.getValues() { //not ordered
    println(s)
}

// get a specific name of the `enum`
println(LogOption.getName(LogOption.Lshortfile))
```

### Control flow

* If-else
* For/for-in
* while
* do
* try-catch-finally
* case-in/case-is

```swift
// if-else
let a= 10, b = 5
if (a > b) { // '()' is optional, but '{}' is a must
    println("a > b")
}
elseif a == b { // could also use 'elsif'
    println("a = b")
}
else {
    println("a < b")
}

// for
i = 9
for { // forever loop
    i = i + 2
    if (i > 20) { break }
    println('i = {i}')
}


i = 0
for (i = 0; i < 5; i++) {  // c-like for, '()' is a must
    if (i > 4) { break }
    if (i == 2) { continue }
    println('i is {i}')
}


for i in range(10) {
    println('i = {i}')
}

a = [1,2,3,4]
for i in a where i % 2 != 0 {
    println(i)
}


hs = {"a"=>1, "b"=>2, "c"=>3, "d"=>4, "e"=>5, "f"=>6, "g"=>7}
for k, v in hs where v % 2 == 0 {
    println('{k} : {v}')
}


for i in 1..5 {
    println('i={i}')
}

for item in 10..20 where $_ % 2 == 0 { // $_ is the index
    printf("idx=%d, item=%d\n", $_, item)
}


for c in "m".."a" {
    println('c={c}')
}


for idx, v in "abcd" {
    printf("idx=%d, v=%s\n", idx, v)
}


for idx, v in ["a", "b", "c", "d"] {
    printf("idx=%d, v=%s\n", idx, v)
}

for item in ["a", "b", "c", "d"] where $_ % 2 == 0 { // $_ is the index
    printf("idx=%d, item=%s\n", $_, v)
}


//for loop is an expression, not statement, so it could be assigned to a variable
let plus_one = for i in [1,2,3,4] { i + 1 }
fmt.println(plus_one)

// while
i = 10
while (i>3) {
    i--
    println('i={i}')
}

// do
i = 10
do {
    i--
    if (i==3) { break }
}

// try-catch-finally(only support string type)
let exceptStr = "SUMERROR"
try {
    let th = 1 + 2
    if (th == 3) { throw exceptStr }
}
catch "OTHERERROR" {
    println("Catched OTHERERROR")
}
catch exceptStr {
    println("Catched is SUMERROR")
}
catch {
    println("Catched ALL")
}
finally {
    println("finally running")
}

// case-in/case-is
let testStr = "123"
case testStr in { // in(exact/partial match), is(only exact match)
    "abc", "mno" { println("testStr is 'abc' or 'mno'") }
    "def"        { println("testStr is 'def'") }
    `\d+`        { println("testStr contains digit") }
    else         { println("testStr not matched") }
}

let i = [{"a"=>1, "b"=>2}, 10]
let x = [{"a"=>1, "b"=>2},10]
case i in {
    1, 2 { println("i matched 1, 2") }
    3    { println("i matched 3") }
    x    { println("i matched x") }
    else { println("i not matched anything")}
}

```

### Array

In monkey, you could use [] to initialize an empty array:

```swift
emptyArr = []
emptyArr[3] = 3 //will auto expand the array
println(emptyArr)
```

Array could contain any number of different data types.

```swift
mixedArr = [1, 2.5, "Hello", ["Another", "Array"], {"Name"=>"HHF", "SEX"=>"Male"}]
```

You could use index to access array element.

```swift
println('mixedArr[2]={mixedArr[2]}')
println(["a", "b", "c", "d"][2])
```

Because array is an object, so you could use the object's method to operate on it.

```swift
if ([].empty()) {
    println("array is empty")
}

emptyArr.push("Hello")
println(emptyArr)

//you could also use 'addition' to add an item to an array
emptyArr += 2
println(emptyArr)
```

Array could be iterated using `for` loop

```swift
numArr = [1,3,5,2,4,6,7,8,9]
for item in numArr where item % 2 == 0 {
    println(item)
}

let strArr = ["1","a5","5", "5b","4","cc", "7", "dd", "9"]
for item in strArr where /^\d+/.match(item) {
    println(item)
}

for item in ["a", "b", "c", "d"] where $_ % 2 == 0 {  //$_ is the index
    printf("idx=%d, v=%s\n", $_, item)
}
```

You could also use the builtin `reverse` function to reverse array element:

```swift
let arr = [1,3,5,2,4,6,7,8,9]
println("Source Array =", arr)

revArr = reverse(arr)
println("Reverse Array =", revArr)
```

### String

In monkey, there are three types of `string`:

* Raw string
* Double quoted string(Could not contains newline)
* Single quoted string(Interpolated String)

Raw string literals are character sequences between back quotes, as in `foo`. Within the quotes, any character may appear except back quote.

See below for some examples:

```swift
normalStr = "Hello " + "world!"
println(normalStr)

println("123456"[2])

rawStr = `Welcome to
visit us!`
println(rawStr)

//when you use single quoted string, and want variable to be interpolated,
//you just put the variable into '{}'. see below:
str = "Hello world"
println('str={str}') //output: "Hello world"
str[6]="W"
println('str={str}') //output: "Hello World"

```

In monkey, strings are utf8-encoded, you could use utf-8 encoded name as a variable name.

```swift
三 = 3
五 = 5
println(三 + 五) //output : 8
```

strings are also object, so you could use some of the methods provided by `strings` module.

```swift
upperStr = "hello world".upper()
println(upperStr) //output : HELLO WORLD
```

string could also be iterated:

```swift
for idx, v in "abcd" {
    printf("idx=%d, v=%s\n", idx, v)
}

for v in "Hello World" {
    printf("idx=%d, v=%s\n", $_, v) //$_ is the index
}
```

You could concatenate an object to a string:

```swift
joinedStr = "Hello " + "World"
joinedStr += "!"
println(joinedStr)
```

You could also can use the builtin `reverse` function to reverse character of the string:

```swift
let str = "Hello world!"
println("Source Str =", str)
revStr = reverse(str)
println("Reverse str =", revStr)
```

### Hash

In monkey, you could use {} to initialize an empty hash:

```swift
emptyHash = {}
emptyHash["key1"] = "value1"
println(emptyHash)
```

Hash's key could be string, int, boolean:

```swift
hashObj = {
    12     => "twelve",
    true   => 1,
    "Name" => "HHF"
}
println(hashObj)
```

You could use '+' or '-' to add or remove an item from a hash:

```swift
hashObj += {"key1" => "value1"}
hashObj += {"key2" => "value2"}
hashObj += {5 => "five"}
hashObj -= "key2"
hashObj -= 5
println(hash)
```

In monkey, Hash is also an object, so you could use them to operate on hash object:

```swift

hashObj.push(15, "fifteen") //first parameter is the key, second is the value
hashObj.pop(15)

keys = hashObj.keys()
println(keys)

values = hashObj.values()
println(values)
```

You could also use the builtin `reverse` function to reverse hash's key and value:

```swift
let hs = {"key1"=>12, "key2"=>"HHF", "key3"=>false}
println("Source Hash =", hs)
revHash = reverse(hs)
println("Reverse Hash =", revHash)
```

## Standard input/output/error

There are three predefined object for representing standard input, standard output, standard error.
They are `stdin`, `stdout`, `stderr`.

```swift
stdout.writeLine("Hello world")
//same as above
fmt.fprintf(stdout, "Hello world\n")

print("Please type your name:")
name = stdin.read(1024)  //read up to 1024 bytes from stdin
println("Your name is " + name)
```

## Error Handling of standard library

When a standard library function returns `nil` or `false`, you can use the return value's message() function for the error message:

```swift
file = newFile(filename, "r")
if (file == nil) {
    println("opening ", filename, "for reading failed, error:", file.message())
}


let ret = http.listenAndServe("127.0.0.1:9090")
if (ret == false) {
    println("listenAndServe failed, error:", ret.message())
}

```

Maybe you are curious about why `nil` or `false` have message() function? Because in monkey, `nil` and `false`
both are objects, so they have method to operate on it.

## About `defer` keyword

A defer statement defers the execution of a function until the surrounding function returns.

The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns.

```swift
let add  =  fn(x,y){
    defer println("I'm defer1")
    println("I'm in add")
    defer println("I'm defer2")
    return x + y
}
println(add(2,2))
```

The result is as below:

```sh
I'm in add
I'm defer2
I'm defer1
4
```

```swift
file = newFile(filename, "r")
if (file == nil) {
    println("opening ", filename, "for reading failed, error:", file.message())
    return false
}
defer file.close()
//do other file related stuff, and not need to worry about the file close.
//when any file operation error occurs, it will close the file before it returns.

```

## Concatenation of different types

In monkey, you could concatenate of different types. See below for examples:

```swift
// Number plus assignment
num = 10
num += 10 + 15.6
num += 20
println(num)

// String plus assignment
str = "Hello "
str += "world! "
str += [1, 2, 3]
println(str)

// Array plus assignment
arr = []
arr += 1
arr += 10.5
arr += [1, 2, 3]
arr += {"key"=>"value"}
println(arr)

// Array compare
arr1 = [1, 10.5, [1, 2, 3], {"key" => "value"}]
println(arr1)
if arr == arr1 { //support ARRAY compare
    println("arr1 = arr")
} else {
    println("arr1 != arr")
}

// Hash assignment("+=", "-=")
hash = {}
hash += {"key1" => "value1"}
hash += {"key2" => "value2"}
hash += {5 => "five"}
println(hash)
hash -= "key2"
hash -= 5
println(hash)
```

## grep and map

The `grep` and `map` operators are just like perl's `grep` and `map`.

The grep operator takes a list of values and a "testing expression." For each item in the list of values,
the item is placed temporarily into the $_ variable, and the testing expression is evaluated. If the
expression results in a true value, the item is considered selected.

The map operator has a very similar syntax to the grep operator and shares a lot of the same operational steps.
For example, items from a list of values are temporarily placed into $_ one at a time. However,
the testing expression becomes a mapping expression.

```swift
let sourceArr = [2,4,6,8,10,12]

let m = grep  $_ > 5, sourceArr
println('m is {m}')

let cp = map $_ * 2 , sourceArr
println('cp is {cp}')

//a little bit more complex example
let fields = {
                "animal"   => "dog",
                "building" => "house",
                "colour"   => "red",
                "fruit"    => "apple"
             }
let pattern = `animal|fruit`
// =~(match), !~(unmatch)
let values = map { fields[$_] } grep { $_ =~ pattern } fields.keys()
println(values)
```

## Function

Function in monkey is a first-class object. This means the language supports passing functions as arguments to
other functions, returning them as the values from other functions, and assigning them to variables or storing
them in data structures.

Function also could have default parameters and variadic parameters.

```swift
//define a function
let add = fn() { [5,6] }
let n = [1, 2] + [3, 4] + add()
println(n)


let complex = {
   "add" => fn(x, y) { return fn(z) {x + y + z } }, //function with closure
   "sub" => fn(x, y) { x - y },
   "other" => [1,2,3,4]
}
println(complex["add"](1, 2)(3))
println(complex["sub"](10, 2))
println(complex["other"][2])


let warr = [1+1, 3, fn(x) { x + 1}(2),"abc","def"]
println(warr)


println("\nfor i in 5..1 where i > 2 :")
for i in fn(x){ x+1 }(4)..fn(x){ x+1 }(0) where i > 2 {
  if (i == 3) { continue }
  println('i={i}')
}


// default parameter and variadic parameters
add = fn (x, y=5, z=7, args...) {
    w = x + y + z
    for i in args {
        w += i
    }
    return w
}

w = add(2,3,4,5,6,7)
println(w)
```

## Pipe Operator

The pipe operator, inspired by [Elixir](https://elixir-lang.org/).
And thanks for the project [Aria](https://github.com/fadion/aria), I got the idea and some code from this project.

See below for examples:

```swift
# Test pipe operator(|>)
x = ["hello", "world"] |> strings.join(" ") |> strings.upper() |> strings.lower() |> strings.title()
printf("x=<%s>\n", x)

let add = fn(x,y) { return x + y }
let pow = fn(x) { return x ** 2}
let subtract = fn(x) { return x - 1}

let mm = add(1,2) |> pow() |> subtract()
printf("mm=%d\n", mm)
```

## Spawn and channel

You can use `spawn` to create a new thread, and `chan` to communicate with the thread.

```swift
let aChan = chan()
spawn fn() {
    let message = aChan.recv()
    println('channel received message=<{message}>')
}()

//send message to thread
aChan.send("Hello Channel!")
```

## Standard module introduction

In monkey, there are some standard modules provided for you. e.g. json, sql, sort, fmt, os, logger, time, flag, net, http, etc...
This is a brief introduction of some of the monkey standard modules, don't expect it to be thorough.
If you are curious, please see the source code.

### fmt module

```swift
let i = 108, f = 25.383, b=true, s = "Hello, world",
    aArr = [1, 2, 3, 4, "a", "b"],
    aHash = { "key1" => 1, "key2" => 2, "key3" => "abc"}

// Use '%v (value)' to print variable value, '%_' to print the variable's type
fmt.printf("i=[%05d, %X], b=[%t], f=[%.5f], s=[%-15s], aArr=%v, aHash=%v\n", i, i, b, f, s, aArr, aHash)
fmt.printf("i=[%_], b=[%t], f=[%f], aArr=%_, aHash=%_, s=[%s] \n", i, b, f, aArr, aHash, s)

sp = fmt.sprintf("i=[%05d, %X], b=[%t], f=[%.5f], s=[%-15s]\n", i, i, b, f, s)
fmt.printf("sp=%s", sp)

fmt.fprintf(stdout, "Hello %s\n", "world")
```

#### time module

```swift
t1 = newTime()
format = t1.strftime("%F %R")
println(t1.toStr(format))
Epoch = t1.toEpoch()
println(Epoch)

t2 = t1.fromEpoch(Epoch)
println(t2.toStr(format))
```

#### logger module

```swift
#Log to stdout
log = newLogger(stdout, "LOGGER-", logger.LSTDFLAGS | logger.LMICROSECONDS)

log.printf("Hello, %s\n", "logger")
fmt.printf("Logger: flags =<%d>, prefix=<%s>\n", log.flags(), log.prefix())

#Log to file
file = newFile("./logger.log", "a+")
log.setOutput(file)
for i in 1..5 {
    log.printf("This is <%d>\n", i)
}
```

#### flag module(for handling of command line options)

```swift
let verV = flag.bool("version", false, "0.1")
let ageV = flag.int("age", 40, "an int")
let heightV = flag.float("height", 120.5, "a float")
let nameV = flag.string("name", "HuangHaiFeng", "a string")
let hobbiesV = flag.string("hobbies", "1,2,3", "a comma-delimited string")
flag.parse()

println("verV = ", verV)
println("ageV = ", ageV)
println("heightV = ", heightV)
println("nameV = ", nameV)
println("hobbies = ", hobbiesV.split(","))

if (flag.isSet("age")) {
    println("age is set")
} else {
    println("age is not set")
}
```

#### json module( for json marshal & unmarshal)

```swift
let hsJson = {"key1" => 10,
              "key2" => "Hello Json %s %s Module",
              "key3" => 15.8912,
              "key4" => [1,2,3.5, "Hello"],
              "key5" => true,
              "key6" => {"subkey1"=>12, "subkey2"=>"Json"},
              "key7" => fn(x,y){x+y}(1,2)
}
let hashStr = json.marshal(hsJson) //same as `json.toJson(hsJson)`
println(json.indent(hashStr, "  "))

let hsJson1 = json.unmarshal(hashStr)
println(hsJson1)


let arrJson = [1,2.3,"HHF",[],{ "key" =>10, "key1" =>11}]
let arrStr = json.marshal(arrJson)
println(json.indent(arrStr))
let arr1Json = json.unmarshal(arrStr)  //same as `json.fromJson(arrStr)`
println(arr1Json)
```

#### net module

```swift
//A simple tcp client
let conn = dialTCP("tcp", "127.0.0.1:9090")
if (conn == nil) {
    println("dailTCP failed, error:", conn.message())
    os.exit(1)
}

let n = conn.write("Hello server, I'm client")
if (n == nil) {
    println("conn write failed, error:", n.message())
    os.exit(1)
}

let ret = conn.close()
if (ret == false) {
    println("Server close failed, error:", ret.message())
}

//A simple tcp server
let ln = listenTCP("tcp", ":9090")
for {
    let conn = ln.acceptTCP()
    if (conn == nil) {
        println(conn.message())
    } else {
        printf("Accepted client, Address=%s\n", conn.addr())
    }
    spawn fn(conn) { //spawn a thread to handle the connection
        println(conn.read())
    }(conn)

} //end for

let ret = ln.close()
if (ret == false) {
    println("Server close failed, error:", ret.message())
}
```

#### linq module

```swift
let mm = [1,2,3,4,5,6,7,8,9,10]
println('before mm={mm}')

result = linq.from(mm).where(fn(x) {
    x % 2 == 0
}).select(fn(x) {
    x = x + 2
}).toSlice()
println('after result={result}')

result = linq.from(mm).where(fn(x) {
    x % 2 == 0
}).select(fn(x) {
    x = x + 2
}).last()
println('after result={result}')

let sortArr = [1,2,3,4,5,6,7,8,9,10]
result = linq.from(sortArr).sort(fn(x,y){
    return x > y
})
println('[1,2,3,4,5,6,7,8,9,10] sort(x>y)={result}')

result = linq.from(sortArr).sort(fn(x,y){
    return x < y
})
println('[1,2,3,4,5,6,7,8,9,10] sort(x<y)={result}')

thenByDescendingArr = [
    {"Owner" => "Google",    "Name" => "Chrome"},
    {"Owner" => "Microsoft", "Name" => "Windows"},
    {"Owner" => "Google",    "Name" => "GMail"},
    {"Owner" => "Microsoft", "Name" => "VisualStudio"},
    {"Owner" => "Google",    "Name" => "GMail"},
    {"Owner" => "Microsoft", "Name" => "XBox"},
    {"Owner" => "Google",    "Name" => "GMail"},
    {"Owner" => "Google",    "Name" => "AppEngine"},
    {"Owner" => "Intel",     "Name" => "ParallelStudio"},
    {"Owner" => "Intel",     "Name" => "VTune"},
    {"Owner" => "Microsoft", "Name" => "Office"},
    {"Owner" => "Intel",     "Name" => "Edison"},
    {"Owner" => "Google",    "Name" => "GMail"},
    {"Owner" => "Microsoft", "Name" => "PowerShell"},
    {"Owner" => "Google",    "Name" => "GMail"},
    {"Owner" => "Google",    "Name" => "GDrive"}
]

result = linq.from(thenByDescendingArr).orderBy(fn(x) {
    return x["Owner"]
}).thenByDescending(fn(x){
    return x["Name"]
}).toOrderedSlice()    //Note: You need to use toOrderedSlice

//use json.indent() for formatting the output
let thenByDescendingArrStr = json.marshal(result)
println(json.indent(thenByDescendingArrStr, "  "))

//test 'selectManyByIndexed'
println()
let selectManyByIndexedArr1 = [[1, 2, 3], [4, 5, 6, 7]]
result = linq.from(selectManyByIndexedArr1).selectManyByIndexed(
fn(idx, x){
    if idx == 0 { return linq.from([10, 20, 30]) }
    return linq.from(x)
}, fn(x,y){
    return x + 1
})
println('[[1, 2, 3], [4, 5, 6, 7]] selectManyByIndexed() = {result}')

let selectManyByIndexedArr2 = ["st", "ng"]
result = linq.from(selectManyByIndexedArr2).selectManyByIndexed(
fn(idx,x){
    if idx == 0 { return linq.from(x + "r") }
    return linq.from("i" + x)
},fn(x,y){
    return x + "_"
})
println('["st", "ng"] selectManyByIndexed() = {result}')
```

#### csv module

```swift
//test csv reader
let r = newCsvReader("./examples/test.csv")
if r == nil {
    printf("newCsv returns err, message:%s\n", r.message())
}

r.setOptions({"Comma"=>";", "Comment"=>"#"})

ra = r.readAll()
if (ra == nil) {
    printf("readAll returns err, message:%s\n", ra.message())
}

for line in ra {
    println(line)
    for record in line {
        println("	", record)
    }
}

//test csv writer
let ofile = newFile("./examples/demo.csv", "a+")
let w = newCsvWriter(ofile)
w.setOptions({"Comma"=>"	"})
w.write(["1", "2", "3"])
w.writeAll([["4", "5", "6"],["7", "8", "9"],["10", "11", "12"]])
w.flush()
```

#### template module

The `template` module contains 'text' and 'html' template handling.

Use `newText(...)` or `parseTextFiles(...)` to create a new 'text' template.

Use `newHtml(...)` or `parseHtmlFiles(...)` to create a new 'html' template.

```swift
arr = [
    { "key" => "key1", "value" => "value1" },
    { "key" => "key2", "value" => "value2" },
    { "key" => "key3", "value" => "value3" }
]

//use parseTextFiles(), write to a string
template.parseTextFiles("./examples/looping.tmpl").execute(resultValue, arr)
println('{resultValue}')

//use parseTextFiles(), write to a file
file = newFile("./examples/outTemplate.log", "a+")
template.parseTextFiles("./examples/looping.tmpl").execute(file, arr)

//use parse()
//Note here: we need to use "{{-" and "-}}" to remove the newline from the output
template.newText("array").parse(`Looping
{{- range . }}
        key={{ .key }}, value={{ .value -}}
{{- end }}
`).execute(resultValue, arr)
println('{resultValue}')
```

#### sql module

The `sql` module provides a lower abstraction layer for working with database.

It should correctly handle database null values, though not throughly tested.

For testing `sql` module, you need to do following:

1. Download sql driver source.

2. Include the package in 'sql.go' like below:

```go
    _ "github.com/mattn/go-sqlite3"
```

3. Recompile monkey source.

Below is a complete source of the `examples/db.my`:

```swift
let dbOp = fn() {
    os.remove("./foo.db") //delete `foo.db` file
    let db = dbOpen("sqlite3", "./foo.db")
    if (db == nil) {
        println("DB open failed, error:", db.message())
        return false
    }
    defer db.close()
    let sqlStmt = `create table foo (id integer not null primary key, name text);delete from foo;`
    let exec_ret = db.exec(sqlStmt)
    if (exec_ret == nil) {
        println("DB exec failed! error:", exec_ret.message())
        return false
    }

    let tx = db.begin()
    if (tx == nil) {
        println("db.Begin failed!, error:", tx.message())
        return false
    }

    let stmt = tx.prepare(`insert into foo(id, name) values(?, ?)`)
    if (stmt == nil) {
        println("tx.Prepare failed!, error:", stmt.message())
        return false
    }

    defer stmt.close()
    let i = 0
    for (i = 0; i < 105; i++) {
        let name = "您好" + i
        if (i>100) {
            //insert `null` value. There are five predefined values:INT_NULL,FLOAT_NULL,STRING_NULL,BOOL_NULL,TIME_NULL.
            let rs = stmt.exec(i, sql.STRING_NULL)
        } else {
            let rs = stmt.exec(i, name)
        }

        if (rs == nil) {
            println("statement exec failed, error:", rs.message())
            return false
        }
    } //end for

    tx.commit()

    let id = 0, name = ""
    let rows = db.query("select id, name from foo")
    if (rows == nil) {
        println("db queue failed, error:", rows.message())
        return false
    }
    defer rows.close()
    while (rows.next()) {
        rows.scan(id, name)
        if (name.valid()) { //check if it's `null`
            println(id, "|", name)
        } else {
            println(id, "|", "null")
        }
    }
    return true
}

let ret = dbOp()
if (ret == nil) {
    os.exit(1)
}

os.exit()
```

## About regular expression

In monkey, regard to regular expression, you could use:

* Regular expression literal
* 'regexp' module
* '=~' and '!~' operators(like perl's)

```swift
//Use regular expression literal( /pattern/.match(str) )
let regex = /\d+\t/.match("abc 123	mnj")
if (regex) { println("regex matched using regular expression literal") }

//Use 'regexp' module
if regexp.compile(`\d+\t`).match("abc 123	mnj") {
    println("regex matched using 'regexp' module")
}

//Use '=~'(str =~ pattern)
if "abc 123	mnj" =~ `\d+\t` {
    println("regex matched using '=~'")
}else {
    println("regex not matched using '=~'")
}

```

```sh
Note: For detailed explanation of 'Regular Expression' pattern matching, you could see golang's regexp module for reference.
```

## Useful Utilities

Included has some useful utilities like `formatter` and `highlighter`.

The formatter utility can format the monkey language.
The highlighter utility can highlight the monkey language to console or html.

You could also combine the two utilities:

```sh
./fmt xx.my | ./highlight  //output to console(console highlight not support windows)
```

## Syntax Highlight

Currently there are three kind of syntax highlight for editors:

1. vim

    [vim](misc/vim)

2. emeditor

    [emeditor](misc/emeditor)

3. notepad++

    [notepad++](misc/notepad%2B%2B)

4. Visual Studio Code

    [VSC](misc/vscode)

5. Sublime Text 3

    [Sublime Text 3](misc/SublimeText3)

## Futual Plans

There are some other things i plan to do:

* Improve the Standard Library with more functions.
* Write more tests!
* Improve this document with more explanation of the language.
* Rewrite the demo program for better understanding of the language.
* Rewrite the 'include' module logic.
* ~~Add support for if-elseif-else expression~~.

## License

MIT
