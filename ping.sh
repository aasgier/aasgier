#!/bin/sh

## CONFIGURATION

# Set if the script runs on the primary or secondary RPi.
primary=false

# Set IP address to ping.
ip='45.32.184.111'

# Set interval in seconds at which to ping.
interval=1

# Set time to wait till an RPi is declared "dead".
wait=1

# Set script to run.
script='/bin/true'

## EXECUTE

ping="$(ping -q -c 1 -W $wait '45.32.184.121')"

while true; do
	if $primary; then
		echo 'Running on: primary'
		$script
		echo "Script executed with exit status: $?"
	else
		echo 'Running on: secondary'
		if echo "$ping" | grep "1 received" > /dev/null; then
			# TODO: Add exit status here
			echo "Script executed on primary with exit status: "
		else
			$script
			echo "Primary is down, script executed on secondary with exit status: $?"
		fi
	fi

	sleep $interval
done
