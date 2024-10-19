#!/bin/bash

read_pid() {
  file_path="$1"

if [[ -f "$file_path" ]]; then
  pid=$(grep -o -E '[0-9]+(\.[0-9]+)?' "$file_path" | head -n 1)
  if [[ -n "$first_number" ]]; then
    echo "First number in file: $first_number"
  else
    echo "There are no numbers in the file."
  fi
else
  echo "File does not exist."
fi

}

start() {
  file_path="${MONITORING_DIR}/${FILE_NAME}"
  if [[ -f "$file_path" ]]; then
    pid=$(grep -o -E '[0-9]+(\.[0-9]+)?' "$file_path" | head -n 1)
    if [[ -n "$pid" ]]; then
      echo "monitoring is already running: process ${pid}"
      exit 0
    fi
  fi

  touch "$file_path"

  bash "$MONITOR_FILE" "$MONITORING_DIR" &

  echo "$!" > "$file_path"
}

stop() {
  file_path="${MONITORING_DIR}/${FILE_NAME}"
  if [[ -f "$file_path" ]]; then
    pid=$(grep -o -E '[0-9]+(\.[0-9]+)?' "$file_path" | head -n 1)
    if [[ -n "$pid" ]]; then
      kill "$pid"
      rm "$file_path"
      exit 0
    fi
  fi

  echo "monitoring is not running"
}

status() {
  file_path="${MONITORING_DIR}/${FILE_NAME}"
  if [[ -f "$file_path" ]]; then
    pid=$(grep -o -E '[0-9]+(\.[0-9]+)?' "$file_path" | head -n 1)
    if [[ -n "$pid" ]]; then
      echo "monitoring is running: process ${pid}"
      exit 0
    fi
  fi

  echo "monitoring is not running"
}

if [ "$#" -ne 1 ]; then
    echo "Error: Expected exactly one argument"
    exit 1
fi

MONITORING_DIR="${HOME}/custom_disk_monitoring"
mkdir -p "$MONITORING_DIR"
FILE_NAME=".monitor_pid.txt"
MONITOR_FILE="monitor.sh"

case "$1" in
  START)
    start
    ;;
  STOP)
    stop
    ;;
  STATUS)
    status
    ;;
  *)
    echo "Unknown command: try smth from START, STOP, STATUS"
    exit 1
    ;;
esac

exit 0

