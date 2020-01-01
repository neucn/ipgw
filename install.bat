@echo off
if exist ipgw.exe (
    if exist %APPDATA%\ipgw\ipgw.exe del %APPDATA%\ipgw\ipgw.exe
    if not exist %APPDATA%\ipgw md %APPDATA%\ipgw
    echo Moving ipgw.exe
    copy ipgw.exe %APPDATA%\ipgw\
    echo Editing Environment Variable
    setx /m ipgw "%APPDATA%\ipgw"
    start SystemPropertiesAdvanced.exe
    echo Done
) else (
   echo Can't find ipgw.exe in this dir
   echo Fail
)
pause