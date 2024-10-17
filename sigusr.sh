#!/bin/bash


PROCESS_NAME=myapp
RESTART=false
SCRIPT_PATH=""
PID=$(pgrep -f "$PROCESS_NAME")


function show_help() {

  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  -r            Restart the server (send SIGUSR1)"
  echo "  -s <script>   Path to the SQL script to run (send SIGUSR2)"
  echo "  -h            Show this help message"
  exit 0
}


while getopts ":rhs:" opt; do
  case $opt in
    r) RESTART=true ;;
    s) SCRIPT_PATH="$OPTARG" ;;
    h) show_help ;;
    \?) echo "Invalid option: -$OPTARG" >&2; exit 1 ;;
    :) echo "Option -$OPTARG requires an argument." >&2; exit 1;;
  esac
done


function find_pid() {
  pgrep -f "$PROCESS_NAME"
}

function send_signal() {
  local signal=$1
  local pid=$2

  if [ -n "$pid" ]; then
    kill -$signal "$pid"
    echo "Sent $signal to process $pid"
  else
    echo "Process not found"
    exit 1
  fi
}



PID=$(find_pid)

if $RESTART; then
  send_signal SIGUSR1 "$PID"
  exit 0
fi

if [ -n "$SCRIPT_PATH" ]; then
  send_signal SIGUSR2 "$PID"

  sleep 2

  if [ ! -f "$SCRIPT_PATH" ]; then
    echo "SQL script not found at path: $SCRIPT_PATH"
    cat "Empty" | tee /proc/$PID/fd/0
    exit 0

  fi

  cat "$SCRIPT_PATH" | tee /proc/$PID/fd/0
  echo "SQL script sent to the process"
  exit 0
fi

show_help
