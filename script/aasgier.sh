#!/bin/fish

# Simulate the time it takes to check the sensors.
<<<<<<< HEAD
sleep 1
=======
sleep 4
>>>>>>> 20d9cd492398dd707e01d659a7f433e6963b97c7

# Stuff the python scrip should print.
# I use toml syntax for this.
echo 'waterlevel = '(random 0 100)
<<<<<<< HEAD
echo 'vibrate = '(echo -e 'true\nfalse' | shuf -n 1)
=======
echo 'windspeed = '(random 0 100)
echo 'barrierstatus = "'(echo -e 'open\nclosed' | shuf -n 1)'"'
>>>>>>> 20d9cd492398dd707e01d659a7f433e6963b97c7
