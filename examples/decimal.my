d1 = decimal.fromString("123.45678901234567")
d2 = decimal.fromFloat(3)
decimal.setDivisionPrecision(50)
fmt.println("123.45678901234567/3 = ", d1.div(d2))
fmt.println(d1.div(d2))

fmt.println(decimal.fromString("123.456").trunc(2))

//convert string to decimal
d3=decimal("123.45678901234567")
fmt.println(d3)
fmt.println("123.45678901234567/3 = ", d3.div(d2))

//convert decimal to string
str = str(d3.div(d2))
fmt.println(str)

//convert decimal to int, it will only return the int part of the decimal
i = int(d3)
fmt.println(i)

//convert decimal to float, it may lost precision.
f = float(d3)
fmt.println(f)

//Be careful with below sentence, it will output '125.12512312312312',
//not '125.1251231231231231231231231231231231231', because you use
// 'fromFloat', the golang runtime will first convert 
//'125.1251231231231231231231231231231231231' to '125.12512312312312'.
//You should use decimal.fromString("125.1251231231231231231231231231231231231")
//instead.
println(decimal.fromFloat(125.1251231231231231231231231231231231231)) //not what you expected