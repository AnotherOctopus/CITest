#!/bin/bash
# this takes the raspiconfig functions and hardcodes them to be turned
# on at run. WARNING, WE NEED TO KEEP UP TO DATE WITH RASPICONFIG
# raspiconfig can be found here
# https://github.com/asb/raspi-config/blob/master/raspi-config

set_camera() {
  # Stop if /boot is not a mountpoint
  if ! mountpoint -q /boot; then
    return 1
  fi

  [ -e $CONFIG ] || touch $CONFIG

  if [ "$1" -eq 0 ]; then # disable camera
    set_config_var start_x 0 $CONFIG
    sed $CONFIG -i -e "s/^startx/#startx/"
    sed $CONFIG -i -e "s/^start_file/#start_file/"
    sed $CONFIG -i -e "s/^fixup_file/#fixup_file/"
  else # enable camera
    set_config_var start_x 1 $CONFIG
    CUR_GPU_MEM=$(get_config_var gpu_mem $CONFIG)
    if [ -z "$CUR_GPU_MEM" ] || [ "$CUR_GPU_MEM" -lt 128 ]; then
      set_config_var gpu_mem 128 $CONFIG
    fi
    sed $CONFIG -i -e "s/^startx/#startx/"
    sed $CONFIG -i -e "s/^fixup_file/#fixup_file/"
  fi
}

if [ ! -e /boot/start_x.elf ]; then
	echo "Your firmware appears to be out of date (no start_x.elf). Please update" 20 60 2
	return 1
fi
set_camera 1;
