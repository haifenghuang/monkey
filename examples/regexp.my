# test regexp literal(regexp)
let ww = /\d+/.match("12") //regular expression
println(ww)  # true
if (ww) {
    println("matched!")
}
else {
    println("Not matched!")
}

let srcStr = "abc 12 def 34"
let destStr = /\d+/.replace(srcStr, "TT")
println(destStr)  # "abc TT def TT"

let regex = /\d+\t/.match("abc 123	mnj")
if (regex) {println("regex matched")}

let arr = / /.split("ba na za")
if (len(arr) > 0) {
	println("/ /.split('ba na za')[0]=", arr[0])
}
else {
	println("Not splitted")
}

let orgSrc = "Tomato is fruit"
let patternA = `fruit`
# =~(match), !~(unmatch)
if orgSrc =~ patternA {
	println("'Tomato is fruit' contains 'fruit'")
}else {
	println("'Tomato is fruit' not contains 'fruit'")
}

let patternB = `animal`
if orgSrc != patternB {
	println("'Tomato is fruit' not contains 'animal'")
}else {
	println("'Tomato is fruit' contains 'animal'")
}

//test ReplaceAllStringFunc()
lowerText = "antonio carlos jobim"
//expr = regexp.mustCompile("\\b\\w")
expr = regexp.mustCompile("[a-z]")
upperText = expr.replaceAllStringFunc(lowerText, fn(x){
	return strings.upper(x)
})
println('replaceAllStringFunc, (antonio carlos jobim), upper=({upperText})')


//test findAllStringSubmatch()
re = regexp.compile(`(.)at`) // want to know what is in front of 'at'
res = regexp.findAllStringSubmatch("The cat sat on the mat.", -1)
fmt.printf("findAllStringSubmatch, %v\n", res) //[[cat c] [sat s] [mat m]]



if regexp.mustCompile(`Hello`).matchString("Hello Regular Expression.") == true {
	fmt.printf("Match\n")
} else {
	fmt.printf("No match\n")
}

