#!/bin/sh

# Download files from ftp://wirelessftp.fcc.gov/pub/uls/complete/
# Files to be downloaded are identified in config.json

file='config.json'

cat "$file" | jq -r 'keys[]' | 

while IFS= read -r file; do
  if [ ! -d "wirelessftp.fcc.gov/pub/uls/complete/$file" ] 
  then
    wget -m ftp://wirelessftp.fcc.gov/pub/uls/complete/$file.zip
    for i in wirelessftp.fcc.gov/pub/uls/complete/*.zip; do unzip "$i" -d "${i%%.zip}"; done
    rm -rf wirelessftp.fcc.gov/pub/uls/complete/*.zip
  fi
done


exit 