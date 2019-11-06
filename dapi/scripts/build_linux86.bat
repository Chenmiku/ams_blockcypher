cd ..
set arch=386
set dir=release
set filename=dapi_linux86
del %dir%\%filename%
bash.exe -c -l "export GOARCH=%arch% && go build -o %dir%/%filename%"