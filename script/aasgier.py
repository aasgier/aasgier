#!/usr/bin/python3

from random import randint

# Get sensor information.
# TODO: Use actual TI functions for this.
waterLevel = randint(50,60)
vibrate = randint(0,1)
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

# Convert vibrate bool to bool accepted by toml.
if vibrate:
    vibrate = "false"
if not vibrate:
    vibrate = "true"

# Print toml stuff that gets parsed by go program.
print("waterlevel =", waterLevel)
print("windspeed =", windSpeed)
print("vibrate =", vibrate)
#if closed:
#    print("barriers = ","closed")
#else:
#    print("barriers = ", "open")
