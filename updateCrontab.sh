#!/bin/bash

export PYTHONIOENCODING=utf8

content=$(curl http://worklog-dev.odds.team/api/v1/reminder/setting)

# get date from setting
date=$(echo $content | python -c \
       "import sys, json; \
        print json.load(sys.stdin)['setting']['date']")
# get hour from setting
hour=$(echo $content | python -c \
       "import sys, json; \
       print json.load(sys.stdin)['setting']['time'].split(':')[0]")
# get minute from setting
min=$(echo $content | python -c \
       "import sys, json; \
       print json.load(sys.stdin)['setting']['time'].split(':')[1]")

# set default date
if [ -z $date ] 
then 
    date=25 
fi

# set default hour
if [ -z $hour ] 
then 
    hour=23 
fi

# set default min
if [ -z $min ] 
then 
    min=59 
fi

# create crontab reminder
crontab -l > reminder

# setting up crontab
echo "$min $hour $date * * ./callApi.sh" > reminder
#echo "*/1 * * * * ./callApi.sh" > reminder

# run crontab
crontab reminder

# remove temporary crontab file
rm reminder

    