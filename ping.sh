#!/bin/sh

## CONFIGURATION

# Set if the script runs on the primary or secondary RPi.
primary=true

# Set IP address of RPi to ping.
ip='45.32.184.111'

# Set interval in seconds at which to ping.
interval=1

# Set time to wait till an RPi is declared "dead".
wait=1

# Set script to run.
script='/bin/true'

# File to redirect exit status to.
file='/tmp/exit'


## EXECUTE

ping=$(ping -q -c 1 -W "$wait" "$ip")

# This functions pinrt a string to stderr if $1 (the exit status) is not 0.
print () {
	if test "$1" -eq 0; then
		echo "$2"
	else
		echo "$2" >&2
	fi
}

while true; do
	if "$primary"; then
		echo 'Running on: primary'

		"$script"

		status="$?"
		print "$status" "Primary is down, script executed on secondary with exit status: $status"
		echo "$status" > "$file"
	else
		echo 'Running on: secondary'

		if echo "$ping" | grep '1 received' >'/dev/null'; then
			status=$(ssh "TODO@$ip" cat "$file")
			print "$status" "Primary is down, script executed on secondary with exit status: $status"
		else
			"$script"

			status="$?"
			print "$status" "Primary is down, script executed on secondary with exit status: $status"
			echo "$status" > "$file"
		fi
	fi

	echo
	sleep $interval
done
