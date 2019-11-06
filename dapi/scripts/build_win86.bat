cd ..
set GOARCH=386
set dir=release
set filename=dapi_win86.exe
del %dir%\%filename%
go build -o %dir%\%filename%