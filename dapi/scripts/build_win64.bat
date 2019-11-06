cd ..
set GOARCH=amd64
set dir=release
set filename=dapi_win64.exe
del %dir%\%filename%
go build -o %dir%\%filename%