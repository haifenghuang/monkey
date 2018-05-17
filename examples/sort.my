let romanNumeralDict = {
  1000: "M",
  900 : "CM",
  500 : "D",
  400 : "CD",
  100 : "C",
  90  : "XC",
  50  : "L",
  40  : "XL",
  10  : "X",
  9   : "IX",
  5   : "V",
  4   : "IV",
  1   : "I",
}

keys = []
for k, v in romanNumeralDict {
    keys += k
}

keysNum = sort.sortInts(keys)
for k in keysNum {
    fmt.printf("%v\t: %v\n", k, romanNumeralDict[k])
}


let stringsDict = {
  "1000": "M",
  "900" : "CM",
  "500" : "D",
  "400" : "CD",
  "100" : "C",
  "90"  : "XC",
  "50"  : "L",
  "40"  : "XL",
  "10"  : "X",
  "9"   : "IX",
  "5"   : "V",
  "4"   : "IV",
  "1"   : "I",
}

println("<===========SORT ASENDING=============>")
keys = []
for k, v in stringsDict {
    keys += k
}

//default is ascending
keyAsend = sort.sortStrings(keys)
for k in keyAsend {
    fmt.printf("%v\t: %v\n", k, stringsDict[k])
}


println("<===========SORT DESCENDING=============>")
// There are four predefined variable for sorting order
//    sort.Ascending(default)
//    sort.Descending
//    sort.CaseInsensitiveAscending
//    sort.CaseInsensitiveDescending
keyDescend = sort.sortStrings(keys, sort.Descending)
for k in keyDescend {
    fmt.printf("%v\t: %v\n", k, stringsDict[k])
}