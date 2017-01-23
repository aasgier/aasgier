#!/bin/fish

# Simulate the time it takes to check the sensors.
#sleep 1


# Stuff the python scrip should print.
# I use toml syntax for this.
echo 'waterlevel = '(random 50 60)
echo 'vibrate = '(echo -e 'true\nfalse' | shuf -n 1)
