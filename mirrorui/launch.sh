#!/bin/bash
npm start &
blah=$$
echo $blah
sleep 5
npm run electron
pkill $blah
