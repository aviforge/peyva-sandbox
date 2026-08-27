package content

// RunnerScript is the ready-made script that starts and stops the processes
// that make up peyva.
//
// The book hands this over rather than asking an assistant to write it, for two
// reasons. It is boilerplate: process supervision teaches nothing this book is
// about, and every reader would pay tokens for a slightly different version of
// the same sixty lines. And asking for it read as a contradiction, because the
// preamble on that prompt says to build in the reader's chosen language and the
// runner is the one thing that cannot be.
//
// The scripts know nothing about the reader's language. They set environment
// variables and run commands the reader fills in at the top, which is the
// whole contract:
//
//   - PEYVA_PORT    the port this process listens on. Every process reads it.
//   - PEYVA_VAULT   the port the Vault listens on. The copies read it.
//   - PEYVA_PEERS   the copies' ports, comma separated. The proxy reads it.
//   - PEYVA_PRIMARY the primary Vault's port. The replica reads it.
//   - PEYVA_WARDEN  the Warden's port. The Vaults and the copies read it.
//
// Two of the commands start blank. The replica and the Warden do not exist
// when the script is handed over, and a command that starts nothing is how the
// script grows a process when a later chapter builds one.
type RunnerScript struct {
	// SystemID matches System.ID.
	SystemID string
	// Path is where the reader saves it.
	Path string
	// Content is the script.
	Content string
}

// RunnerChapter is where the script is handed over: the chapter that first runs
// more than one process. Which system it is written for was chosen back in
// chapter 0, beside the language.
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
// running behind it. The copies sit above it, the stores below it.
const runBash = `#!/usr/bin/env bash
# peyva/run.sh - start, inspect and stop the processes that make up peyva.
#   ./run.sh start 3   the Vault, three copies, and the proxy in front
#   ./run.sh status    what is alive
#   ./run.sh stop      kill everything peyva started

# Job control, on purpose. It puts each background job in its own process group,
# which is what makes stop able to kill a process and the language runtime under
# it together. Without it they all share this script's group and stop reaches
# only the outermost shell, leaving the port held.
set -muo pipefail
cd "$(dirname "$0")/.."

# Set these to how your project starts. Every process reads PEYVA_PORT. The
# copies read PEYVA_VAULT, the proxy reads PEYVA_PEERS. Two start blank and are
# filled in by the chapter that builds them: the replica also reads
# PEYVA_PRIMARY, and the Vaults and the copies read PEYVA_WARDEN once it exists.
START_VAULT="go run ./peyva/vault"
START_REPLICA=""
START_WARDEN=""
START_COPY="go run ./peyva/gateway"
START_PROXY="go run ./peyva/proxy"

VAULT_PORT=9300
REPLICA_PORT=9301
WARDEN_PORT=9302
PROXY_PORT=9310
FIRST_COPY_PORT=9311
RUNDIR="peyva/.run"
PIDFILE="$RUNDIR/groups"
PORTFILE="$RUNDIR/ports"

# Tagging is a read loop rather than sed because sed buffers its output when it
# is not writing to a terminal, and a reader watching three copies start would
# see nothing at all until one of them stopped.
prefix() {
  while IFS= read -r line; do echo "[$1] $line"; done
}

# Pure bash, so status needs neither lsof nor netstat, neither of which is
# reliably present.
port_open() {
  (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null && exec 3<&- && return 0
  return 1
}

# launch <tag> <port> <command> [VAR=value ...]
# Every process gets PEYVA_PORT and the two addresses it might need. A process
# that does not read a variable is not harmed by it being set.
launch() {
  name="$1"; port="$2"; cmd="$3"; shift 3
  ( env PEYVA_PORT="$port" PEYVA_VAULT="$VAULT_PORT" PEYVA_WARDEN="$WARDEN_PORT" "$@" $cmd 2>&1 | prefix "$name" ) &
  echo $! >> "$PIDFILE"
  echo "$port" >> "$PORTFILE"
}

start() {
  n="${1:-3}"
  if [ -s "$PIDFILE" ]; then echo "peyva is already running. Run stop first."; exit 1; fi
  mkdir -p "$RUNDIR"
  : > "$PIDFILE"
  : > "$PORTFILE"

  # Stores first, so a copy's first request has somewhere to go.
  launch vault "$VAULT_PORT" "$START_VAULT"
  [ -n "$START_REPLICA" ] && launch replica "$REPLICA_PORT" "$START_REPLICA" PEYVA_PRIMARY="$VAULT_PORT"
  [ -n "$START_WARDEN" ] && launch warden "$WARDEN_PORT" "$START_WARDEN"
  sleep 1

  peers=""
  for i in $(seq 0 $((n - 1))); do
    peers="${peers:+$peers,}$((FIRST_COPY_PORT + i))"
  done
  for i in $(seq 0 $((n - 1))); do
    port=$((FIRST_COPY_PORT + i))
    launch "copy $port" "$port" "$START_COPY"
  done

  launch proxy "$PROXY_PORT" "$START_PROXY" PEYVA_PEERS="$peers"

  echo "vault on $VAULT_PORT, $n copies on $peers, proxy on $PROXY_PORT. Ctrl+C stops them."
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
# has taken some of the processes with it.
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

const runPowerShell = `# peyva/run.ps1 - start, inspect and stop the processes that make up peyva.
#   .\run.ps1 start 3   the Vault, three copies, and the proxy in front
#   .\run.ps1 status    what is alive
#   .\run.ps1 stop      kill everything peyva started
param([string]$Command = "", [int]$Count = 3)
Set-Location (Join-Path $PSScriptRoot "..")

# Set these to how your project starts. Every process reads PEYVA_PORT. The
# copies read PEYVA_VAULT, the proxy reads PEYVA_PEERS. Two start blank and are
# filled in by the chapter that builds them: the replica also reads
# PEYVA_PRIMARY, and the Vaults and the copies read PEYVA_WARDEN once it exists.
$StartVault   = "go run ./peyva/vault"
$StartReplica = ""
$StartWarden  = ""
$StartCopy    = "go run ./peyva/gateway"
$StartProxy   = "go run ./peyva/proxy"

$VaultPort     = 9300
$ReplicaPort   = 9301
$WardenPort    = 9302
$ProxyPort     = 9310
$FirstCopyPort = 9311
$PidFile       = "peyva/.run/pids"

# Each process gets its own shell with the variables set for it alone, so they
# cannot leak between processes or outlive the run. Single quotes: $env:... must
# reach the child as text to assign, not be expanded to this shell's own value
# before it gets there.
function Start-One([int]$port, [string]$cmd, [string]$extra = "") {
  $assign = '$env:PEYVA_PORT=' + $port + '; $env:PEYVA_VAULT=' + $VaultPort + '; $env:PEYVA_WARDEN=' + $WardenPort + '; '
  $p = Start-Process -PassThru -NoNewWindow powershell -ArgumentList @(
    "-NoProfile", "-Command", ($assign + $extra + $cmd))
  return $p.Id
}

function Start-Peyva([int]$n) {
  if (Test-Path $PidFile) { Write-Host "peyva is already running. Run stop first."; return }
  New-Item -ItemType Directory -Force -Path "peyva/.run" | Out-Null
  $ids = @()

  # Stores first, so a copy's first request has somewhere to go.
  $ids += Start-One $VaultPort $StartVault
  if ($StartReplica -ne "") { $ids += Start-One $ReplicaPort $StartReplica ('$env:PEYVA_PRIMARY=' + $VaultPort + '; ') }
  if ($StartWarden -ne "")  { $ids += Start-One $WardenPort $StartWarden }
  Start-Sleep -Seconds 1

  $peers = @()
  for ($i = 0; $i -lt $n; $i++) {
    $port = $FirstCopyPort + $i
    $ids += Start-One $port $StartCopy
    $peers += $port
  }

  $list = $peers -join ","
  $ids += Start-One $ProxyPort $StartProxy ('$env:PEYVA_PEERS=' + "'$list'" + '; ')

  $ids | Set-Content -Encoding utf8 $PidFile
  Write-Host "vault on $VaultPort, $n copies on $list, proxy on $ProxyPort. Run stop when you are done."
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
// It tracks processes by port rather than by process id: a batch file has no
// clean way to capture the id of something it starts, and the port is what has
// to be free before peyva can start again anyway.
const runBatch = `@echo off
REM peyva\run.bat - start, inspect and stop the processes that make up peyva.
REM   run.bat start 3   the Vault, three copies, and the proxy in front
REM   run.bat status    what is alive
REM   run.bat stop      kill everything peyva started
setlocal enabledelayedexpansion
cd /d "%~dp0.."

REM Set these to how your project starts. Every process reads PEYVA_PORT. The
REM copies read PEYVA_VAULT, the proxy reads PEYVA_PEERS. Two start blank and
REM are filled in by the chapter that builds them: the replica also reads
REM PEYVA_PRIMARY, and the Vaults and the copies read PEYVA_WARDEN once it
REM exists.
set "START_VAULT=go run ./peyva/vault"
set "START_REPLICA="
set "START_WARDEN="
set "START_COPY=go run ./peyva/gateway"
set "START_PROXY=go run ./peyva/proxy"

set "VAULT_PORT=9300"
set "REPLICA_PORT=9301"
set "WARDEN_PORT=9302"
set "PROXY_PORT=9310"
set "FIRST_COPY_PORT=9311"
set "RUNDIR=peyva\.run"
set "PORTFILE=%RUNDIR%\ports"
set "NETFILE=%RUNDIR%\net.tmp"
set "COMMON=set PEYVA_VAULT=%VAULT_PORT%&& set PEYVA_WARDEN=%WARDEN_PORT%"

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

REM Stores first, so a copy's first request has somewhere to go.
echo %VAULT_PORT%>> "%PORTFILE%"
start "peyva vault" /min cmd /c "set PEYVA_PORT=%VAULT_PORT%&& %COMMON%&& %START_VAULT%"
if defined START_REPLICA (
  echo %REPLICA_PORT%>> "%PORTFILE%"
  start "peyva replica" /min cmd /c "set PEYVA_PORT=%REPLICA_PORT%&& set PEYVA_PRIMARY=%VAULT_PORT%&& %COMMON%&& %START_REPLICA%"
)
if defined START_WARDEN (
  echo %WARDEN_PORT%>> "%PORTFILE%"
  start "peyva warden" /min cmd /c "set PEYVA_PORT=%WARDEN_PORT%&& %COMMON%&& %START_WARDEN%"
)
timeout /t 1 /nobreak >nul

set "PEERS="
set /a LAST=%N%-1
for /l %%i in (0,1,!LAST!) do (
  set /a PORT=%FIRST_COPY_PORT%+%%i
  if defined PEERS (set "PEERS=!PEERS!,!PORT!") else (set "PEERS=!PORT!")
  echo !PORT!>> "%PORTFILE%"
  start "peyva copy !PORT!" /min cmd /c "set PEYVA_PORT=!PORT!&& %COMMON%&& %START_COPY%"
)

echo %PROXY_PORT%>> "%PORTFILE%"
start "peyva proxy" /min cmd /c "set PEYVA_PORT=%PROXY_PORT%&& set PEYVA_PEERS=!PEERS!&& %COMMON%&& %START_PROXY%"

echo vault on %VAULT_PORT%, %N% copies on !PEERS!, proxy on %PROXY_PORT%. Run stop when you are done.
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
REM has already taken some of the processes with it. /T takes the language
REM runtime under each window with it.
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
