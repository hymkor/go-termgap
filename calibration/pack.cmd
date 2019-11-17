setlocal
set GOARCH=386
go generate
go fmt
go build -ldflags "-s -w"
for %%I in (%CD%) do set NAME=%%~nI
zip %NAME%-%DATE:/=%-windows-386.zip %NAME%.exe
endlocal
