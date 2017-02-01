#!/usr/bin/python3

import csv
import RPi.GPIO as GPIO
import os
import time

buttonA = 5
buttonB = 6
buttonC = 13
ledYellow = 12
ledGreen = 16
motor = 18
buttonD = 19
klopfSensor = 21
ledRed = 26
gateGOTO = True
averageWaterLevel = []
averageWaterLevels = 0
data = {}
vibrationLevel = 'veilig'
windSpeed = 0

GPIO.setmode(GPIO.BCM)
GPIO.setwarnings(False)
GPIO.setup(motor, GPIO.OUT)
GPIO.setup(ledRed, GPIO.OUT)
GPIO.setup(ledYellow, GPIO.OUT)
GPIO.setup(ledGreen, GPIO.OUT)
GPIO.setup(buttonA, GPIO.IN)
GPIO.setup(buttonB, GPIO.IN)
GPIO.setup(buttonC, GPIO.IN)
GPIO.setup(buttonD, GPIO.IN)
GPIO.setup(klopfSensor, GPIO.IN)

def operateGate(gateGOTO):
        global p
        p = GPIO.PWM(18, 50)
        p.start(7.5)
        if (gateGOTO == False):
                p.ChangeDutyCycle(2.5)
                gateGOTO = True
                print('Kering geopend.')
        elif (gateGOTO == True):
                p.ChangeDutyCycle(12.5)
                gateGOTO = False
                print('Kering gesloten.')
        return gateGOTO

def getSensorData():
        global data
        global vibrationLevel
        global windSpeed
        global gateGOTO
        getWaterLevel()
        data['Waterniveau'] = averageWaterLevel
        print('Waterniveau gemeten.')
        getVibration()
        data['Grondvibratie'] = vibrationLevel
        getWindSpeed()
        data['Windsnelheid'] = windSpeed
        print(data)
        operateGate(gateGOTO)
        return data

def readADC(adcnum, clockpin, mosipin, misopin, cspin):
        global averageWaterLevel
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
    global averageWaterLevel
    global averageWaterLevels
    averageWaterLevels = []
    # change these as desired - they're the pins connected from the
    # SPI port on the ADC to the Cobbler
    SPICLK = 11
    SPIMISO = 9
    SPIMOSI = 10
    SPICS = 8

    # set up the SPI interface pins
    GPIO.setup(SPIMOSI, GPIO.OUT)
    GPIO.setup(SPIMISO, GPIO.IN)
    GPIO.setup(SPICLK, GPIO.OUT)
    GPIO.setup(SPICS, GPIO.OUT)

    # 10k trim pot connected to adc #0
    potentiometer_adc = 0;
    readCount = 0

    while readCount <= 4:
            # read the analog pin
            set_volume = readADC(potentiometer_adc, SPICLK, SPIMOSI, SPIMISO, SPICS)
            averageWaterLevels.append(int(set_volume))
            # hang out and do nothing for a half second
            readCount += 1
            time.sleep(0.5)
    averageWaterLevel = sum(averageWaterLevels)/len(averageWaterLevels)
    return averageWaterLevel

def getVibration():
        global vibrationLevel
        if (GPIO.input(klopfSensor) == False):
                vibrationLevel = 'Onveilig'
        return vibrationLevel

def getWindSpeed():
        global windSpeed
        print('Kies windsnelheid:')
        while True:
                if (GPIO.input(buttonD) == True):
                        GPIO.output(26, False)
                        GPIO.output(12, False)
                        GPIO.output(16, False)
                        #print("knop D doet het goed")
                        with open("windSpeed.csv", "w", newline="") as csvfile:
                            writer = csv.writer(csvfile, delimiter=",")
                            writer.writerow(["windSpeed"] + [0])
                            csvfile.close()
                            break
                elif (GPIO.input(buttonC) == True):
                        #print("knop C doet het goed")
                        GPIO.output(26, False)
                        GPIO.output(12, False)
                        GPIO.output(16, True)
                        with open("windSpeed.csv", "w", newline="") as csvfile:
                            writer = csv.writer(csvfile, delimiter=",")
                            writer.writerow(["windSpeed"] + [3])
                            csvfile.close()
                        windSpeed = 3
                        break
                elif (GPIO.input(buttonB) == True):
                        #print("knop B doet het goed")
                        GPIO.output(26, False)
                        GPIO.output(16, False)
                        GPIO.output(12, True)
                        with open("windSpeed.csv", "w", newline="") as csvfile:
                            writer = csv.writer(csvfile, delimiter=",")
                            writer.writerow(["windSpeed"] + [6])
                            csvfile.close()
                        windSpeed = 6
                        break
                elif (GPIO.input(buttonA) == True):
                        #print("knop A doet het goed")
                        GPIO.output(26, True)
                        GPIO.output(12, False)
                        GPIO.output(16, False)
                        with open("windSpeed.csv", "w", newline="") as csvfile:
                            writer = csv.writer(csvfile, delimiter=",")
                            writer.writerow(["windSpeed"] + [9])
                            csvfile.close()
                        windSpeed = 9
                        break
        return windSpeed
