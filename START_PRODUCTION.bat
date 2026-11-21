@echo off
title Pulse Production Server & Cloudflare Tunnel
cd /d "%~dp0"

echo ========================================================
echo  PULSE WORKSPACE - PRODUCTION HOST & CLOUDFLARE TUNNEL
echo  Target URL: https://project.kedukanaksara.com
echo ========================================================

:: 1. Jalankan Cloudflare Tunnel di latar belakang
start /b "" "C:\Users\agrap\Documents\GitHub\ChatbotBaru\.tools\cloudflared.exe" tunnel run 67775b37-cc35-49c5-9cf2-7135fabea95a

:: 2. Jalankan Node Production Gateway (Gateway Port 5080 + Go Backend Port 5081)
node production_server.js
