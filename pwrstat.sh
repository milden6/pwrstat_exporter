#!/bin/bash

# Start daemon
/etc/init.d/pwrstatd start

# Disable UPS power fail and low battery, because we need only status from it
pwrstat -pwrfail -active off -shutdown off
pwrstat -lowbatt -runtime 0 -capacity 0 -active off -shutdown off