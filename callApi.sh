#!/bin/bash

# call api
result=$(curl http://worklog-dev.odds.team/api/v1/reminder/send)

#convert value from response to success | failed
[[ $result = true ]] && state="Success" || state="Failed"

# write result to file
echo "[$(date "+%Y-%m-%d %H:%M:%S")] --> Send result:" $state >> reminder_send_logs.txt