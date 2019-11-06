cd ..
set arch=amd64
set dir=release
set filename=dapi_linux64
del %dir%\%filename%
bash.exe -c -l "export GOARCH=%arch% && go build -o %dir%/%filename%"