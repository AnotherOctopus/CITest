#!/bin/bash
docker build . -t xxcore
docker tag xxcore anotheroctopus/xxcore
docker login -u anotheroctopus -p anotheroctopus/xxcore
docker push anotheroctopus/xxcore
ssh pi@128.46.156.193 'docker pull anotheroctopus/xxcore'
ssh pi@128.46.156.193 'docker run anotheroctopus/xxcore'
