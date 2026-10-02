@echo off
set MSXGL_PATH=..\MSXgl
if not exist out ( md out )
node %MSXGL_PATH%\engine\script\js\build.js %*
