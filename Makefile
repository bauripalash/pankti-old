# Build Pankti
#

all:
	go build

wasm:
	GOOS=js GOARCH=wasm go build -o wasm/res/pankti.wasm wasm/panktiWasm.go

pyserver:
	cd wasm/res/ && python -m http.server 8099

win:
	GOOS=windows GOARCH=amd64 go build -o dist/pankti_x86-64.exe --tags noide

win32:
	GOOS=windows GOARCH=386 go build -o dist/pankti_x86.exe --tags noide 

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/pankti_x86-64 --tags noide

linux32:
	GOOS=linux GOARCH=386 go build -o dist/pankti_x86 --tags noide 
 
