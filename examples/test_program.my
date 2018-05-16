let times_two = fn(n) {
  '{n} * 2 = {n * 2 }'
}

let do_loop = fn() {
  let x = 0  
  do {
    x = x + 1
    println(times_two(x))
    if (x == 10) {
      println("x is 10!")
      break
    }
  }
}

let y = 1; let z = 5;
do {
  if (y == 5) {
    println("y is 5")
    do_loop()
    break
  }
  y = y + 1
}

z = "a b c d e f".split(" ").map(fn(x) {
  println(x.upper())
  x.upper()
})
println(z)

let w = 5
#test continue
do {
    w = w + 1
    if (w > 15) {
        break
    }
    if (w == 10) {
      println("w is 10!, continue")
      continue
    }
}
println('w is {w}')

while (w > 0) {
    w = w - 1
    println('w is {w}')
}

let H={"hello":"hello", "world":"world", "123":"456"}
println(H["hello"])

# test grep & map
let mm = [1,2,3,4,5,6,7,8,9,10]
println('before mm={mm}')
mm = mm.grep(fn(x) {  //test for grep 
  x % 2 == 0
}).map(fn(x) {  # test for map
  x = x * 2
})
println('after mm={mm}')

let k =  10.5
k = k + 8
println('k={k}')

let pp = 3 * "hello " #test string with '*'
println('pp={pp}')

let kk = 5.1 + "hello"
print("kk is ", kk, "!  pp is ", pp)
println("")

let subs = "Hello World".substr(6, 10)
println(subs)

# test for loop
let y22 = 0
for x in [1,3,5,7,9] {
    y22 = 2 * x
    if (x >= 7) { break }
    println('x is {x}, y22 is {y22}')
}

for key, value in H {
    println("key=[", key, "] value=[",value, "]")
}

let i = 5
for i in range(10) {
    println('i = {i}')
}

println('outer i = {i}')

i = 3
for (i; i < 8; i = i+1) {
    if (i == 6) { continue }
    println('i = <{i}>')
}

i = 9
for {
    i =  i + 2
    if (i > 20) { break }
    println('i = ({i})')
}

println("END OF TEST")
