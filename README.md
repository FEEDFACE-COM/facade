# FACADE by FEEDFACE.COM
    
## Examples

#### system status: `top`
```
# raspi #
facade render term -shape wave
# client #
facade exec -host raspi term -w 80 -h 25 top -1

```


#### network traffic: `tshark`
```
# raspi #
facade render words -shape field -n 32 -life 4 -mark 1 -shuffle 
# client #
sudo tshark -i wlan0 -l -T fields -e ip.src \
| nc raspi 4045

```


#### logfiles: `tail`
```
# raspi #
facade render lines -shape disk -w 150 -h 12
# client #
tail -f /var/log/nginx/access.log \
| nc raspi 4045


```

#### trace route: `mtr`
```
# raspi #
facade serve term -shape vortex
# client #
facade exec term -w 120 -h 16 sudo mtr -m 10 --displaymode 2 wikipedia.org


```


#### wall time: `date`
```
# raspi #
facade render chars -shape moebius -w 64 -speed .5 -font spacemono 
# client #
while true; do date +"%Y-%m-%dT%H:%M:%S%z"; sleep 1; done \
| nc raspi 4045

```

```
# raspi #
facade serve lines -shape wave -h 2 -w 10 -down -font ocraext -zoom .8
# client #
while true; do date "+%Y-%m-%d"; sleep 1; date "+ %H:%M:%S"; sleep 1; done \
| nc raspi 4045

```


#### shell sharing: `bash`

```
# raspi #
facade render term -mask=f
# client #
facade exec -host raspi term -w 80 -h 25 bash

```


#### text adventures: `frotz`

```
# raspi #
facade render term -shape slate -zoom .75
# client #
facade exec -host raspi term -w 110 -h 30 frotz /path/to/hitchhikers_guide.z5
```


#### some man pages are quite pretty: `man`
```
# raspi #
facade render lines -w 50 -shape crawl
# client #
MANWIDTH=50 MANPAGER=cat man ssh \
 | while read line; do echo "$line"; sleep .9; done \
 | nc raspi 4045
```


#### internetworking specifications in plain text format: `rfc`
```
# raspi #
facade render lines -w 72 -shape rows
# client #
curl -L https://tools.ietf.org/rfc/rfc792.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| nc raspi 4045
```


#### your favourite hacking zine articles: `PHRACK`
```
# raspi #
facade render lines -w 80 -shape roll
# client #
curl -sL http://phrack.org/archives/tgz/phrack49.tar.gz \
| tar xfz /dev/stdin --to-stdout ./14.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| nc raspi 4045
```


#### nudes older than the `<IMG>` tag: `asciipr0n`
```
# raspi #
facade serve lines -w 80 -shape slate 
# client #
curl -sL https://www.asciipr0n.com/pr0n/pinups/pinup00.txt \
| while read -r line; do echo "$line"; sleep .5; done \
| nc raspi 4045

```


#### demo scene release notes: `.nfo`
```
# raspi #
facade render lines -w=80 -vert=wave -mask=mask -font adore64
# client #
curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| nc raspi 4045
```


#### w
```
# raspi #
facade render term -shape slate
# client #
facade exec term -w 50 -h 20 curl parrot.live
```

## Author

If you enjoy FACADE, tell us how you are using it at <facade@feedface.com>!

