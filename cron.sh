#!/usr/bin/env bash

# add this to /etc/cron.hourly/mbo.sh
# #!/usr/bin/env bash
# MBO_USERNAME=username MBO_PASSWORD=password MBO_STUDIO=12375 MBO_TIME="7:00 am" /path/to/cron.sh >> /var/log/mbo_cron.log 2>&1

set -e

mbo login -u $MBO_USERNAME -p $MBO_PASSWORD -studio $MBO_STUDIO

# See if there's a class open to be registered for
today=`date +"%m/%d/%Y"`
class=`mbo ls -date $today -open | grep "$MBO_TIME" | head -n1`
#class="Wed Feb 19 7:00 am    (11 Reserved, 7 Open) 119      CrossFit All Levels    Jenny Werba                   1 hour"

if [ -z "$class" ]; then
  echo "No class open for registration"
  exit 0
fi

# Parse class ID and date
echo $class
id=`echo $class | sed -rn 's/(.*)\)[ ]+([0-9]+)(.*)/\2/p'`
if [ -z "$id" ]; then
  echo Error parsing id: $class
  exit 1
fi
date=`echo $class | sed -rn 's/.{3} (.{3} [0-9+]{2})(.*)\)[ ]+([0-9]+)(.*)/\1/p'`
if [ -z "$date" ]; then
  echo Error parsing date: $class
  exit 1
fi
date=`date --date="$date" +"%m/%d/%Y"`

# register
echo Found class $id on $date, attempting to register
mbo register -date $date -id $id
