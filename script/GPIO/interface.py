#!/usr/bin/python3

import RPi.GPIO as GPIO
import os
import time

buttonA = 5
buttonB = 6
buttonC = 13
buttonD = 19

ledYellow = 12
ledGreen = 16
ledRed = 26

motor = 18
klopfSensor = 21
gateGOTO = True

GPIO.setwarnings(False)
GPIO.setmode(GPIO.BCM)

GPIO.setup(buttonA, GPIO.IN)
GPIO.setup(buttonB, GPIO.IN)
GPIO.setup(buttonC, GPIO.IN)
GPIO.setup(buttonD, GPIO.IN)

GPIO.setup(ledYellow, GPIO.OUT)
GPIO.setup(ledGreen, GPIO.OUT)
GPIO.setup(ledRed, GPIO.OUT)

GPIO.setup(motor, GPIO.OUT)
GPIO.setup(klopfSensor, GPIO.IN)

def operateGate(gateGOTO):
    global p
    p = GPIO.PWM(18, 50)
    p.start(10)

    # True is close!
    # TODO: Use while loop to make it one fluid motion.
    if (gateGOTO == False):
        p.ChangeDutyCycle(8)
    elif (gateGOTO == True):
        p.ChangeDutyCycle(16)

    return gateGOTO

def getSensorData():
    #global gateGOTO
    data = {}

    data['waterLevel'] = getWaterLevel()
    data['vibration'] = getVibration()
    data['windSpeed'] = getWindSpeed()
    #operateGate(gateGOTO)
    return data

def readADC(adcnum, clockpin, mosipin, misopin, cspin):
    if ((adcnum > 7) or (adcnum < 0)):
        return -1
    GPIO.output(cspin, True)

    GPIO.output(clockpin, False)  # start clock low
    GPIO.output(cspin, False)     # bring CS low

    commandout = adcnum
    commandout |= 0x18  # start bit + single-ended bit
    commandout <<= 3    # we only need to send 5 bits here
    for i in range(5):
        if (commandout & 0x80):
            GPIO.output(mosipin, True)
        else:
            GPIO.output(mosipin, False)
        commandout <<= 1
        GPIO.output(clockpin, True)
        GPIO.output(clockpin, False)

    adcout = 0
    # read in one empty bit, one null bit and 10 ADC bits
    for i in range(12):
        GPIO.output(clockpin, True)
        GPIO.output(clockpin, False)
        adcout <<= 1
        if (GPIO.input(misopin)):
            adcout |= 0x1

    GPIO.output(cspin, True)

    adcout >>= 1       # first bit is 'null' so drop it
    return adcout

def getWaterLevel():
    # Change these as desired - they're the pins connected from the
    # SPI port on the ADC to the Cobbler.
    SPICLK = 11
    SPIMISO = 9
    SPIMOSI = 10
    SPICS = 8

    # Set up the SPI interface pins.
    GPIO.setup(SPIMOSI, GPIO.OUT)
    GPIO.setup(SPIMISO, GPIO.IN)
    GPIO.setup(SPICLK, GPIO.OUT)
    GPIO.setup(SPICS, GPIO.OUT)

    # 10k trim pot connected to adc #0.
    potentiometer_adc = 0;

    # Convert water level to something nice (kinda).
    waterLevelRaw = readADC(potentiometer_adc, SPICLK, SPIMOSI, SPIMISO, SPICS)
    waterLevel = round((waterLevelRaw / 1023 * 100) - 20)
    if waterLevel < 1:
        waterLevel = 0

    return waterLevel

def getVibration():
    if (GPIO.input(klopfSensor) == False):
        vibration = 'true'
    else:
        vibration = 'false'

        return vibration

def getWindSpeed():
    windSpeed = 0

    # Change windspeed depending on what button is pressed at the time.
    if (GPIO.input(buttonA) == True):
        GPIO.output(26, True)
        GPIO.output(12, False)
        GPIO.output(16, False)
        windSpeed = 9
    elif (GPIO.input(buttonB) == True):
        GPIO.output(26, False)
        GPIO.output(16, False)
        GPIO.output(12, True)
        windSpeed = 6
    elif (GPIO.input(buttonC) == True):
        GPIO.output(26, False)
        GPIO.output(12, False)
        GPIO.output(16, True)
        windSpeed = 3
    else:
        GPIO.output(26, False)
        GPIO.output(12, False)
        GPIO.output(16, False)

    return windSpeed
