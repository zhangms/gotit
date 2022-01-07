#!/bin/bash

for f in `ls images/*.*g` ; do
  echo "--------BEGIN $f---------"
  cmd="curl -s --user api:DjN0STsWYwbLztxhf0JQjWkh3ZsftD7z --data-binary @./$f -i https://api.tinify.com/shrink"
  resp=`$cmd`
  echo $resp
  url=https:`echo $resp |awk -F"https:" '{print $3}' |awk -F'"' '{print $1}'`
  echo $url
  curl -o ./$f $url
  echo ""
done
