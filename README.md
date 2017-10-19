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
# Integer literals
i1 = 10
i2 = 20_000_000
i3 = 0x80           # hex
i4 = 0b10101        # binary
i5 = 0c127          # octal

# Float literals
f1 = 10.25
f2 = 1.02E3
f3 = 123_456.789_012

# String literals
s1 = "123"
s2 = "Hello world"

# Regular expression literals
r = /\d+/.match("12")
if (r) { prinln("regex matched!") }

# Array literals
a = [1+2, 3, 4, "5", 3]

# Hash literals
h = { "a"=>1, "b"=>2, "c"=>2}

# Nil literal
n = nil

# Boolean literals
t = true
f = false

# Function literals
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

### Type conversion

You can use the builtin `int()`, `float()`, `str()`, `array()` functions for type conversion.

```swift
let i = 0xa
let s = str(i)     # result: "10"
let f = float(i)   # result: 10
let a = array(i)   # result: [10]
```

### Control flow

* If-else
* For/for-in
* while
* do
* try-catch-finally
* case-in/case-is

```swift
# if-else
if (a > b) { # '()' is optional, but '{}' is a must
    println("a > b")
}
else {
    println("a <= b")
}

# for
i = 9
for { # forever loop
    i = i + 2
    if (i > 20) { break }
    println('i = {i}')
}


i = 0
for (i = 0; i < 5; i++) {  # c-like for, '()' is a must
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

for item in 10..20 where $_ % 2 == 0 { # $_ is the index
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

for item in ["a", "b", "c", "d"] where $_ % 2 == 0 { # $_ is the index
    printf("idx=%d, item=%s\n", $_, v)
}

# while
i = 10
while (i>3) {
    i--
    println('i={i}')
}

# do
i = 10
do {
    i--
    if (i==3) { break }
}

# try-catch-finally(only support for string type)
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

# case-in/case-is
let testStr = "123"
case testStr in { #in(exact/partial match), is(only exact match)
    "abc", "mno" { println("testStr is 'abc' or 'mno'") }
    "def"        { println("testStr is 'def'") }
    `\d+`        { println("testStr contains digit") }
    else         { println("testStr not matched") }
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

## License

MIT
