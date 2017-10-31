# Monkey Programming Language

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
* Regular expression support(partially like perls)
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

Monkey supports 7 data types: `String`, `Int`, `Float`, `Bool`, `Array`, `Hash` and `Nil`

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
i2 = 20_000_000
i3 = 0x80           // hex
i4 = 0b10101        // binary
i5 = 0c127          // octal

// Float literals
f1 = 10.25
f2 = 1.02E3
f3 = 123_456.789_012

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
姓="黄"
```

## Reserved keywords

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

Note:

* `where` in expression `for xxx in where` is a reserved keyword, but in other part,
   you could use `where` as a variable, or function name, etc. But it is not recomented.
* `in` in expression `for xxx in`, `case xxx in` is a reserved keyword, but in other part,
   you could use `in` as a variable, or function name, etc. But it is not recomented.
* `is` in expression `case xxx is` is a reserved keyword, but in other part,
   you could use `is` as a variable, or function name, etc. But it is not recomented.

### Type conversion

You can use the builtin `int()`, `float()`, `str()`, `array()` functions for type conversion.

```swift
let i = 0xa
let s = str(i)     // result: "10"
let f = float(i)   // result: 10
let a = array(i)   // result: [10]
```

## `qw`(Quote word) keyword

The `qw` keyword is like perl's `qw` keyword. When you want to use a lot of quoted strings, the `qw` keyword can make it a lot easier for those strings.

```swift
for str in qw<abc, def, ghi, jkl, mno> { //allowed 'qw' pair is '{}', '<>', '()'
  println('str={str}')
}

newArr = qw(1,2,3.5) //array with string values, not number values.
fmt.printf("newArr=%v\n", newArr)
```

## `enum` keyword

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
```

## Standard input/output/error

There are three predefined object for representing standard input, standard output, standard error.
They are `stdin`, `stdout`, `stderr`.

```swift
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

## About concatenation of different types

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

## Grep and map

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

```swift
//fmt module
let i = 108, f = 25.383, b=true, s = "Hello, world",
    aArr = [1, 2, 3, 4, "a", "b"],
    aHash = { "key1" => 1, "key2" => 2, "key3" => "abc"}

// Use '%v (value)' to print variable value, '%_' to print the variable's type
fmt.printf("i=[%05d, %X], b=[%t], f=[%.5f], s=[%-15s], aArr=%v, aHash=%v\n", i, i, b, f, s, aArr, aHash)
fmt.printf("i=[%_], b=[%t], f=[%f], aArr=%_, aHash=%_, s=[%s] \n", i, b, f, aArr, aHash, s)

sp = fmt.sprintf("i=[%05d, %X], b=[%t], f=[%.5f], s=[%-15s]\n", i, i, b, f, s)
fmt.printf("sp=%s", sp)

fmt.fprintf(stdout, "Hello %s\n", "world")


//time module
t1 = newTime()
format = t1.strftime("%F %R")
println(t1.toStr(format))
Epoch = t1.toEpoch()
println(Epoch)

t2 = t1.fromEpoch(Epoch)
println(t2.toStr(format))


//logger module
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


//flag module(for handling of command line options)
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


// json module( for json marshal & unmarshal)
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


//net module
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

## Useful Utilities

Included has some useful utilities like `formatter` and `highlighter`.

The formatter utility can format the monkey language.
The highlighter utility can highlight the monkey language to console or html.

You could also combine the two utilities:

```sh
./fmt xx.my | ./highlight  //output to console(console highlight not support windows)
```

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
