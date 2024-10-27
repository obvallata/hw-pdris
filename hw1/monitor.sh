write_utilization() {
  write_ts="$(date +"%Y-%m-%d %H:%M:%S")"
  df_output="$(df -h)"
  inode_info="$(df -hi)"

  echo "$df_output" | tail -n +2 | while read -r line; do
      filesystem=$(echo "$line" | awk '{print $1}')
      size=$(echo "$line" | awk '{print $2}')
      used=$(echo "$line" | awk '{print $3}')
      available=$(echo "$line" | awk '{print $4}')
      use_percent=$(echo "$line" | awk '{print $5}')

      inode_info=$(echo "$inode_info" | grep "^$filesystem" | awk '{print $2","$3","$4","$5}')

      echo "$write_ts,$filesystem,$size,$used,$available,$use_percent,$inode_info" >> "$FILE_PATH"
  done
}

create_new_file() {
    mkdir -p "$DIRECTORY"

    FILE_PATH="${DIRECTORY}/started_${TIMESTAMP_INIT}_monitor_time_${TIMESTAMP_LAST_FILE}.csv"

    if [ ! -f "$FILE_PATH" ]; then
        touch "$FILE_PATH"
        echo "Timestamp,Filesystem,Size,Used,Available,Use%,Inodes,UsedInodes,FreeInodes,InodesUse%" > "$FILE_PATH"
    fi
}

main() {
  DIRECTORY=$1

  TIMESTAMP_INIT=$(date +%Y%m%d_%H%M%S_%Z)
  TIMESTAMP_LAST_FILE="$TIMESTAMP_INIT"

  create_new_file

   while true; do
      TIMESTAMP_CURRENT=$(date +%Y%m%d_%H%M%S_%Z)

      sec1=$(date -j -f '%Y%m%d_%H%M%S_%Z' "$TIMESTAMP_CURRENT" +'%s')
      sec2=$(date -j -f '%Y%m%d_%H%M%S_%Z' "$TIMESTAMP_LAST_FILE" +'%s')

      difference=$((sec1 - sec2))

      if [ $difference -gt 86400 ]; then
         TIMESTAMP_LAST_FILE="$TIMESTAMP_CURRENT"
         create_new_file
      fi

      write_utilization
      sleep 300
  done
}

if [ $# -ne 1 ]; then
    echo "Expected directory for monitoring files in args"
    exit 1
fi

main "$1"

