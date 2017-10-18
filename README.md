# Monkey Programming Language

**Summary** 

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

```
~ » monkey
Monkey programming language REPL

>>
```

or, to run a program:

```
monkey path/to/file
```
## Useful Utilities
Included has some useful utilities like `formatter` and `highlighter`.

The formatter utility can format the monkey language.
The highlighter utility can highlight the monkey language to console or html.

You could also combine the two utilities:

```
./fmt xx.my | ./highlight  //output to console(console highlight not support windows)
```

## Futual Plans
There are some other things i plan to do:

- Improve the Standard Library with more functions.
- Write more tests!
- Improve this document with more explanation of the language.
- Rewrite the demo program for better understanding of the language.
- Rewrite the 'include' module logic.

## License

MIT
