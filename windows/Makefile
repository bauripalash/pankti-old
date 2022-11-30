all: gui64 gui32

gui64:
	GOOS=windows GOARCH=amd64 go generate
	GOOS=windows GOARCH=amd64 go build -ldflags=all="-H=windowsgui" -o ../dist/panktiw_x86_64.exe 
	rm *.syso

gui32:
	goversioninfo -icon=res/icon.ico -manifest=res/pankti.exe.manifest
	GOOS=windows GOARCH=386 go build -ldflags=all="-H=windowsgui" -o ../dist/panktiw_x86.exe 
	rm *.syso


