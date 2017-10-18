#!/bin/sh
export GOPATH=$(pwd)

# format each go file
#echo "Formatting go file..."
#for file in `find ./src/monkey -name "*.go"`; do
#	echo "    `basename $file`"
#	go fmt $file > /dev/null
#done

echo ""

# run: ./monkey demo.my or ./monkey
echo "Building REPL...(monkey)"
go build -o monkey main.go

# run: ./fmt demo.my
echo "Building Formatter...(fmt)"
go build -o fmt fmt.go

# run:    ./highlight demo.my               (generate: demo.my.html)
#     or  ./fmt demo.my | ./highlight   (generate: output.html)
echo "Building Highlighter...(highlight)"
go build -o highlight highlight.go
