@REM launch this bat under the dir it was in

@echo off
for %%i in (*.proto) do call :setlist %%i

protoc -I=..\..\..\..\ -I=. --go_out=. --go_opt=paths=source_relative %LIST%

echo DONE: %LIST%
set LIST=
exit /b

:setlist
set LIST=%LIST% %~nx1
