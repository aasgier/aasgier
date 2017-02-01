#!/usr/bin/python3

from random import randint
from GPIO.interface import *

# Get sensor information.
data = getSensorData()

# Set initial percentage the water needs to be at for the gate to close.
closep = 90

# Lower closep depending on the circumstances.
if data['vibration']:
    closep -= 20
for i in range(0, data['windSpeed']):
    closep -= 5

if data['waterLevel'] >= closep:
    operateGate(True)
else:
    operateGate(False)

# Print toml stuff that gets parsed by the Go program.
#print("vibration =", data['vibration'])
print("vibration =", str(not bool(randint(0, 3))).lower())
print("waterLevel =", data['waterLevel'])
print("windSpeed =", data['windSpeed'])
