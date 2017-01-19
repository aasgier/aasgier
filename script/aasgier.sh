#!/bin/fish

# Simulate the time it takes to check the sensors.
sleep 4

# Stuff the python scrip should print.
# I use toml syntax for this.
echo 'waterlevel = '(random 0 100)
echo 'windspeed = '(random 0 100)
echo 'barrierstatus = "'(echo -e 'open\nclosed' | shuf -n 1)'"'
