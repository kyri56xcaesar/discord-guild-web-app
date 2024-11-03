#!/bin/bash


PROCESS_NAME=myapp
RESTART=false
SCRIPT_PATH=""
CONFIG_PATH=""


function show_help() {

  echo "Usage: $0 [options]"
  echo "Options:"
  echo "  -p            Specify process name to send the signals"
  echo "  -r            Restart the server (send SIGUSR1)"
  echo "  -s <script>   Path to the SQL script to run (send SIGUSR2)"
  echo "  -c <config>   Path to the new config file for restart"
  echo "  -h            Show this help message"
  exit 0
}


while getopts ":p:rhs:c:" opt; do
  case $opt in
    p) PROCESS_NAME="$OPTARG" ;;
    r) RESTART=true ;;
    s) SCRIPT_PATH="$OPTARG" ;;
    c) CONFIG_PATH="$OPTARG" ;; 
    h) show_help ;;
    \?) echo "Invalid option: -$OPTARG" >&2; exit 1 ;;
    :) echo "Option -$OPTARG requires an argument." >&2; exit 1;;
  esac
done

PID=$(pgrep -f "$PROCESS_NAME")

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


echo "Restart: "$RESTART
echo "Script path: "$SCRIPT_PATH
echo "Config path: "$CONFIG_PATH



if $RESTART; then
  if [ -n "$CONFIG_PATH" ]; then
    export NEW_CONFIG_PATH="$CONFIG_PATH"
  fi
  send_signal SIGUSR1 "$PID"
  exit 0
fi

if [ -n "$SCRIPT_PATH" ]; then

  if [ ! -f "$SCRIPT_PATH" ]; then
    echo "SQL script not found at path: $SCRIPT_PATH"
    cat "Empty" | tee /proc/$PID/fd/0
    exit 0

  fi

  if [ -n "$SCRIPT_PATH" ]; then
    export SQL_SCRIPT_PATH=$SCRIPT_PATH
  fi

  send_signal SIGUSR2 "$PID"


  #cat "$SCRIPT_PATH" | tee /proc/$PID/fd/0
  #echo "SQL script sent to the process"
  exit 0

fi

show_help
