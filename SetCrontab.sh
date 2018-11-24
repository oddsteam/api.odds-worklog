#!/bin/bash

crontab -l > reminder

echo "$1 $2 $3 * * ./callApi.sh" > reminder

crontab reminder

rm reminder