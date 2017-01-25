#!/usr/bin/python3

from random import randint

# Get sensor information.
# TODO: Use actual TI functions for this.
vibrate = randint(0,1)
waterLevel = randint(50,60)
windSpeed = randint(1,12)

# Set initial percentage the water needs to be at for the gate to close.
closep = 100

# Lower closep depending on the circumstances.
if vibrate:
    closep -= 20
for i in range(0, windSpeed):
    closep -= 5
# TODO: Add weather info 'n shit...

if waterLevel >= closep:
    closed = True
else:
    closed = False

# Print toml stuff that gets parsed by the Go program.
print("vibrate =", str(vibrate).lower())
print("waterlevel =", waterLevel)
print("windspeed =", windSpeed)
print("closed =", str(closed).lower())
print("closep =", closep)