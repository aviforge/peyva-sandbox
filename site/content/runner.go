package content

// RunnerScript is the ready-made script that starts and stops copies of peyva.
//
// The book hands this over rather than asking an assistant to write it, for two
// reasons. It is boilerplate: process supervision teaches nothing this book is
// about, and every reader would pay tokens for a slightly different version of
// the same forty lines. And asking for it read as a contradiction, because the
// preamble on that prompt says to build in the reader's chosen language and the
// runner is the one thing that cannot be.
//
// The scripts know nothing about the reader's language. They set two
// environment variables and run a command the reader fills in at the top, which
// is the whole contract: PEYVA_PORT tells a copy which port to listen on, and
// PEYVA_PEERS tells the proxy which copies to route between.
type RunnerScript struct {
	// SystemID matches System.ID.
	SystemID string
	// Path is where the reader saves it.
	Path string
	// Content is the script.
	Content string
}

// RunnerChapter is where the script is handed over: the chapter that first runs
// more than one copy. It is also the chapter that offers the operating system
// picker, because that is the first point at which the answer changes anything.
const RunnerChapter = 10

// RunnerScripts is one entry per operating system. macOS and Linux share a
// script, because nothing it does differs between them.
var RunnerScripts = []RunnerScript{
	{SystemID: "windows", Path: "peyva/run.ps1", Content: runPowerShell},
	{SystemID: "windows-bat", Path: "peyva/run.bat", Content: runBatch},
	{SystemID: "macos", Path: "peyva/run.sh", Content: runBash},
	{SystemID: "linux", Path: "peyva/run.sh", Content: runBash},
}

// RunnerScriptFor returns the script for a system, falling back to the default
// system's rather than returning nothing.
func RunnerScriptFor(systemID string) RunnerScript {
	for _, r := range RunnerScripts {
		if r.SystemID == systemID {
			return r
		}
	}
	for _, r := range RunnerScripts {
		if r.SystemID == DefaultSystem {
			return r
		}
	}
	return RunnerScripts[0]
}

// The proxy takes the port the reader has been using since chapter 2, so every
// URL and bookmark from earlier chapters still works once several copies are
// running behind it. The copies sit above it, out of the way.
const runBash = `#!/usr/bin/env bash
# peyva/run.sh - start, inspect and stop copies of peyva.
#   ./run.sh start 3   three copies behind the proxy
#   ./run.sh status    what is alive
#   ./run.sh stop      kill everything peyva started

# Job control, on purpose. It puts each background job in its own process group,
# which is what makes stop able to kill a copy and the language runtime under it
# together. Without it they all share this script's group and stop reaches only
# the outermost shell, leaving the port held.
set -muo pipefail
cd "$(dirname "$0")/.."

# Set these two to how your project starts. Both read PEYVA_PORT from the
# environment. The proxy also reads PEYVA_PEERS, a comma separated list.
START_COPY="go run ./peyva/gateway"
START_PROXY="go run ./peyva/proxy"

PROXY_PORT=9310
FIRST_COPY_PORT=9311
RUNDIR="peyva/.run"
PIDFILE="$RUNDIR/groups"
PORTFILE="$RUNDIR/ports"

# Tagging is a read loop rather than sed because sed buffers its output when it
# is not writing to a terminal, and a reader watching three copies start would
# see nothing at all until one of them stopped.
tag() {
  while IFS= read -r line; do echo "[$1] $line"; done
}

# Pure bash, so status needs neither lsof nor netstat, neither of which is
# reliably present.
port_open() {
  (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null && exec 3<&- && return 0
  return 1
}

start() {
  n="${1:-3}"
  if [ -s "$PIDFILE" ]; then echo "peyva is already running. Run stop first."; exit 1; fi
  mkdir -p "$RUNDIR"
  : > "$PIDFILE"
  : > "$PORTFILE"

  peers=""
  for i in $(seq 0 $((n - 1))); do
    peers="${peers:+$peers,}$((FIRST_COPY_PORT + i))"
  done

  for i in $(seq 0 $((n - 1))); do
    port=$((FIRST_COPY_PORT + i))
    ( PEYVA_PORT=$port $START_COPY 2>&1 | tag "copy $port" ) &
    echo $! >> "$PIDFILE"
    echo $port >> "$PORTFILE"
  done

  ( PEYVA_PORT=$PROXY_PORT PEYVA_PEERS="$peers" $START_PROXY 2>&1 | tag proxy ) &
  echo $! >> "$PIDFILE"
  echo $PROXY_PORT >> "$PORTFILE"

  echo "$n copies on $peers, proxy on $PROXY_PORT. Ctrl+C stops them."
  trap 'stop; exit 0' INT TERM
  wait
}

status() {
  if [ ! -s "$PORTFILE" ]; then echo "nothing running"; return; fi
  while read -r port; do
    if port_open "$port"; then echo "$port listening"; else echo "$port down"; fi
  done < "$PORTFILE"
}

# A group that has already gone is not an error: stop has to work after a crash
# has taken some of the copies with it.
stop() {
  trap - INT TERM
  if [ ! -s "$PIDFILE" ]; then echo "nothing to stop"; return; fi
  while read -r gid; do kill -TERM "-$gid" 2>/dev/null; done < "$PIDFILE"
  sleep 1
  while read -r gid; do kill -KILL "-$gid" 2>/dev/null; done < "$PIDFILE"
  rm -f "$PIDFILE" "$PORTFILE"
  echo "stopped"
}

case "${1:-}" in
  start) start "${2:-3}" ;;
  status) status ;;
  stop) stop ;;
  *) echo "usage: run.sh start [n] | status | stop"; exit 1 ;;
esac
`

const runPowerShell = `# peyva/run.ps1 - start, inspect and stop copies of peyva.
#   .\run.ps1 start 3   three copies behind the proxy
#   .\run.ps1 status    what is alive
#   .\run.ps1 stop      kill everything peyva started
param([string]$Command = "", [int]$Count = 3)
Set-Location (Join-Path $PSScriptRoot "..")

# Set these two to how your project starts. Both read PEYVA_PORT from the
# environment. The proxy also reads PEYVA_PEERS, a comma separated list.
$StartCopy  = "go run ./peyva/gateway"
$StartProxy = "go run ./peyva/proxy"

$ProxyPort     = 9310
$FirstCopyPort = 9311
$PidFile       = "peyva/.run/pids"

function Start-Peyva([int]$n) {
  if (Test-Path $PidFile) { Write-Host "peyva is already running. Run stop first."; return }
  New-Item -ItemType Directory -Force -Path "peyva/.run" | Out-Null
  $ids = @()
  $peers = @()

  for ($i = 0; $i -lt $n; $i++) {
    $port = $FirstCopyPort + $i
    # Each copy gets its own process with PEYVA_PORT set for it alone, so the
    # variable cannot leak between them or outlive the run.
    # Single quotes: $env:PEYVA_PORT must reach the child as text to assign,
    # not be expanded to this shell's own value before it gets there.
    $cmd = '$env:PEYVA_PORT=' + $port + '; ' + $StartCopy
    $p = Start-Process -PassThru -NoNewWindow powershell -ArgumentList @(
      "-NoProfile", "-Command", $cmd)
    $ids += $p.Id
    $peers += $port
  }

  $list = $peers -join ","
  $cmd = '$env:PEYVA_PORT=' + $ProxyPort + '; $env:PEYVA_PEERS=' + "'$list'" + '; ' + $StartProxy
  $p = Start-Process -PassThru -NoNewWindow powershell -ArgumentList @(
    "-NoProfile", "-Command", $cmd)
  $ids += $p.Id

  $ids | Set-Content -Encoding utf8 $PidFile
  Write-Host "$n copies on $list, proxy on $ProxyPort. Run stop when you are done."
}

function Get-PeyvaStatus {
  if (-not (Test-Path $PidFile)) { Write-Host "nothing running"; return }
  foreach ($id in Get-Content $PidFile) {
    $alive = Get-Process -Id $id -ErrorAction SilentlyContinue
    Write-Host "$id $(if ($alive) { 'alive' } else { 'gone' })"
  }
}

# Kills each recorded process and its children, because the recorded process is
# a shell and the language runtime under it is a child that outlives its parent
# otherwise. A process that has already gone is not an error: stop has to work
# after a crash has taken some of them.
function Stop-Peyva {
  if (-not (Test-Path $PidFile)) { Write-Host "nothing to stop"; return }
  foreach ($id in Get-Content $PidFile) {
    # /T takes the children with it, /F does not ask. Output is discarded
    # because taskkill complains loudly about processes that already exited.
    taskkill /PID $id /T /F 2>$null | Out-Null
  }
  Remove-Item $PidFile -Force
  Write-Host "stopped"
}

switch ($Command) {
  "start"  { Start-Peyva $Count }
  "status" { Get-PeyvaStatus }
  "stop"   { Stop-Peyva }
  default  { Write-Host "usage: run.ps1 start [n] | status | stop" }
}
`

// runBatch exists because a locked-down Windows machine will not run a .ps1.
// It tracks copies by port rather than by process id: a batch file has no clean
// way to capture the id of something it starts, and the port is what has to be
// free before peyva can start again anyway.
const runBatch = `@echo off
REM peyva\run.bat - start, inspect and stop copies of peyva.
REM   run.bat start 3   three copies behind the proxy
REM   run.bat status    what is alive
REM   run.bat stop      kill everything peyva started
setlocal enabledelayedexpansion
cd /d "%~dp0.."

REM Set these two to how your project starts. Both read PEYVA_PORT from the
REM environment. The proxy also reads PEYVA_PEERS, a comma separated list.
set "START_COPY=go run ./peyva/gateway"
set "START_PROXY=go run ./peyva/proxy"

set "PROXY_PORT=9310"
set "FIRST_COPY_PORT=9311"
set "RUNDIR=peyva\.run"
set "PORTFILE=%RUNDIR%\ports"
set "NETFILE=%RUNDIR%\net.tmp"

if /i "%~1"=="start" goto start
if /i "%~1"=="status" goto status
if /i "%~1"=="stop" goto stop
echo usage: run.bat start [n] ^| status ^| stop
exit /b 1

:start
set "N=%~2"
if "%N%"=="" set "N=3"
if exist "%PORTFILE%" echo peyva is already running. Run stop first. & exit /b 1
if not exist "%RUNDIR%" mkdir "%RUNDIR%"
break > "%PORTFILE%"

set "PEERS="
set /a LAST=%N%-1
for /l %%i in (0,1,!LAST!) do (
  set /a PORT=%FIRST_COPY_PORT%+%%i
  if defined PEERS (set "PEERS=!PEERS!,!PORT!") else (set "PEERS=!PORT!")
  echo !PORT!>> "%PORTFILE%"
  start "peyva copy !PORT!" /min cmd /c "set PEYVA_PORT=!PORT!&& %START_COPY%"
)

echo %PROXY_PORT%>> "%PORTFILE%"
start "peyva proxy" /min cmd /c "set PEYVA_PORT=%PROXY_PORT%&& set PEYVA_PEERS=!PEERS!&& %START_PROXY%"

echo %N% copies on !PEERS!, proxy on %PROXY_PORT%. Run stop when you are done.
exit /b 0

REM netstat is written to a file once and read per port. A pipe inside a for /f
REM inside a parenthesised block needs escaping that is easy to get wrong and
REM silently passes garbage to netstat instead of failing.
:status
if not exist "%PORTFILE%" echo nothing running & exit /b 0
netstat -ano -p tcp > "%NETFILE%"
for /f "usebackq delims=" %%p in ("%PORTFILE%") do call :check %%p
del "%NETFILE%" >nul 2>&1
exit /b 0

:check
set "FOUND="
for /f "tokens=5" %%q in ('findstr /r /c:":%~1 .*LISTENING" "%NETFILE%"') do set "FOUND=%%q"
if defined FOUND (echo %~1 listening ^(pid %FOUND%^)) else (echo %~1 down)
exit /b 0

REM A port with nothing on it is not an error: stop has to work after a crash
REM has already taken some of the copies with it. /T takes the language runtime
REM under each window with it.
:stop
if not exist "%PORTFILE%" echo nothing to stop & exit /b 0
netstat -ano -p tcp > "%NETFILE%"
for /f "usebackq delims=" %%p in ("%PORTFILE%") do call :kill %%p
del "%NETFILE%" >nul 2>&1
del "%PORTFILE%" >nul 2>&1
echo stopped
exit /b 0

:kill
for /f "tokens=5" %%q in ('findstr /r /c:":%~1 .*LISTENING" "%NETFILE%"') do taskkill /PID %%q /T /F >nul 2>&1
exit /b 0
`
