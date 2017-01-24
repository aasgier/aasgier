#!/usr/bin/python3

from random import randint

# Get sensor information.
# TODO: Use actual TI functions for this.
vibrate = randint(0,1)
waterLevel = randint(50,60)
windSpeed = randint(0,100)

# Check if the gate shoould be closed or open
# TODO: Use TI function to check if initial closed state should
# be True or False.
if waterLevel > 100 and not vibrate:
    closed = True
elif waterLevel > 75 and vibrate:
    closed = True
else:
    closed = False

# Convert bool to "bool" accepted by toml.
if vibrate:
    vibrate = "true"
else:
    vibrate = "false"
if closed:
    closed = "true"
else:
    closed = "false"

# Print toml stuff that gets parsed by the Go program.
print("vibrate =", vibrate)
print("waterlevel =", waterLevel)
print("windspeed =", windSpeed)
#print("closed = ", closed)
