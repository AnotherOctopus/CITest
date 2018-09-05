#!/bin/bash
# this takes the raspiconfig functions and hardcodes them to be turned
# on at run. WARNING, WE NEED TO KEEP UP TO DATE WITH RASPICONFIG
# raspiconfig can be found here
# https://github.com/asb/raspi-config/blob/master/raspi-config
do_i2c() {
  DEVICE_TREE="yes" # assume not disabled
  DEFAULT=
  if [ -e $CONFIG ] && grep -q "^device_tree=$" $CONFIG; then
    DEVICE_TREE="no"
  fi

  CURRENT_SETTING="off" # assume disabled
  DEFAULT=--defaultno
  if [ -e $CONFIG ] && grep -q -E "^(device_tree_param|dtparam)=([^,]*,)*i2c(_arm)?(=(on|true|yes|1))?(,.*)?$" $CONFIG; then
    CURRENT_SETTING="on"
    DEFAULT=
  fi

  if [ $DEVICE_TREE = "yes" ]; then
    RET=0
    if [ $RET -eq 0 ]; then
      SETTING=on
      STATUS=enabled
    elif [ $RET -eq 1 ]; then
      SETTING=off
      STATUS=disabled
    else
      return 0
    fi
    TENSE=is
    REBOOT=
    if [ $SETTING != $CURRENT_SETTING ]; then
      TENSE="will be"
      REBOOT=" after a reboot"
      ASK_TO_REBOOT=1
    fi
    sed $CONFIG -i -r -e "s/^((device_tree_param|dtparam)=([^,]*,)*i2c(_arm)?)(=[^,]*)?/\1=$SETTING/"
    if ! grep -q -E "^(device_tree_param|dtparam)=([^,]*,)*i2c(_arm)?=[^,]*" $CONFIG; then
      printf "dtparam=i2c_arm=$SETTING\n" >> $CONFIG
    fi
    echo "The ARM I2C interface $TENSE $STATUS$REBOOT" 20 60 1
    if [ $SETTING = "off" ]; then
      return 0
    fi
  fi

  CURRENT_STATUS="yes" # assume not blacklisted
  DEFAULT=
  if [ -e $BLACKLIST ] && grep -q "^blacklist[[:space:]]*i2c[-_]bcm2708" $BLACKLIST; then
    CURRENT_STATUS="no"
    DEFAULT=--defaultno
  fi

  if ! [ -e $BLACKLIST ]; then
    touch $BLACKLIST
  fi

  RET=0
  if [ $RET -eq 0 ]; then
    sed $BLACKLIST -i -e "s/^\(blacklist[[:space:]]*i2c[-_]bcm2708\)/#\1/"
    modprobe i2c-bcm2708
    echo "I2C kernel module will now be loaded by default" 20 60 1
  elif [ $RET -eq 1 ]; then
    sed $BLACKLIST -i -e "s/^#\(blacklist[[:space:]]*i2c[-_]bcm2708\)/\1/"
    if ! grep -q "^blacklist i2c[-_]bcm2708" $BLACKLIST; then
      printf "blacklist i2c-bcm2708\n" >> $BLACKLIST
    fi
    echo "I2C kernel module will no longer be loaded by default" 20 60 1
  else
    return 0
  fi
}
do_i2c
