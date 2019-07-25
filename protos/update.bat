@REM launch this bat with cmd /c update.bat ${workspaceRoot}, you can use it in any proto modules depend on 
@REM the root proto dir (hyperledger.abchain.org/proto)

@echo off
for %%i in (*.proto) do call :setlist %%i

protoc -I=. --go_out=plugins=grpc:%1\src %LIST%

echo DONE: %LIST%

exit /b

:setlist
set LIST=%LIST% %~nx1
