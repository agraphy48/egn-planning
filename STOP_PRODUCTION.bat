@echo off
title Stop Pulse Production Services
echo Menghentikan service Pulse dan Cloudflare Tunnel...

taskkill /f /im node.exe /fi "WINDOWTITLE eq Pulse Production Server*" >nul 2>&1
taskkill /f /im server.exe >nul 2>&1
taskkill /f /im cloudflared.exe >nul 2>&1

echo Service berhasil dihentikan.
pause
