## Setup


## TURN Server
```
$ docker run -it --entrypoint sh -p 3478:3478 -p 49160-49200:49160-49200/udp instrumentisto/coturn
$ /usr/local/bin/docker-entrypoint.sh -n --log-file=stdout --external-ip='$(detect-external-ip)' --min-port=49160 --max-port=49200 --user gamz:gamz
```
