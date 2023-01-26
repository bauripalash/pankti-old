# Build Pankti
#

all: linux_cli

wasm:
	GOOS=js GOARCH=wasm go build -o wasm/res/pankti.wasm wasm/panktiWasm.go

test:
	go test -v -tags noide -cpuprofile cpu.prof -memprofile mem.prof .

bench:
	go test -v -tags noide -cpuprofile cpu.prof -memprofile mem.prof -bench .

pyserver:
	cd wasm/res/ && python -m http.server 8099

win_cli:
	cp windows/versioninfo.json ./
	GOOS=windows GOARCH=amd64 goversioninfo -64 -icon=windows/res/icon.ico -manifest=windows/res/pankti.exe.manifest
	GOOS=windows GOARCH=amd64 go build -o dist/pankti_x86-64.exe --tags noide
	rm versioninfo.json 

win_gui:
	cp windows/versioninfo.json ./
	GOOS=windows GOARCH=amd64 goversioninfo -64 -icon=windows/res/icon.ico -manifest=windows/res/pankti.exe.manifest 
	GOOS=windows GOARCH=amd64 go build -o dist/pankti_x86-64_gui.exe 
	rm versioninfo.json 

win32_cli:
	cp windows/versioninfo.json ./
	GOOS=windows GOARCH=386 goversioninfo -icon=windows/res/icon.ico -manifest=windows/res/pankti.exe.manifest
	GOOS=windows GOARCH=386 go build -o dist/pankti_x86.exe --tags noide 
	rm versioninfo.json

linux_cli:
	GOOS=linux GOARCH=amd64 go build -o dist/pankti_x86-64 --tags noide

linux_gui:
	GOOS=linux GOARCH=amd64 go build -o dist/pankti_x86-64-gui

linux32_cli:
	GOOS=linux GOARCH=386 go build -o dist/pankti_x86 --tags noide 

clean:
	rm -rf dist/*
