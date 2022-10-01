# FACADE by FEEDFACE.COM
    
Pipe text from stdout to FACADE to throw it straight onto the wall of your home/office/hackerspace.

## Examples

#### `top` - system status
```
# on raspi:
facade render term -shape wave

# on client:
facade exec -host raspi term -w 80 -h 25 top -1
```


#### `tshark` - network traffic
```
# on raspi:
facade render words -shape field -n 32 -life 4 -mark 1 -shuffle 

# on client:
alias fcd='nc raspi 4045'
sudo tshark -i wlan0 -l -T fields -e ip.src | fcd
```


#### `tail` - logfiles
```
# on raspi:
facade render lines -shape disk -w 150 -h 12

# on client:
alias fcd='nc raspi 4045'
tail -f /var/log/nginx/access.log | fcd
```

#### `mtr` - trace route
```
# on raspi:
facade render term -shape vortex

# on client:
facade exec term -w 120 -h 16 sudo mtr -m 10 --displaymode 2 wikipedia.org
```


#### `date` - wall time
```
# on raspi:
facade render chars -shape moebius -w 64 -speed .5 -font spacemono

# on client:
alias fcd='nc raspi 4045'
while true; do date +"%Y-%m-%dT%H:%M:%S%z"; sleep 1; done | fcd
```

```
# on raspi:
facade render lines -shape wave -h 2 -w 10 -down -font ocraext -zoom .8

# on client:
alias fcd='nc raspi 4045'
while true; do date "+%Y-%m-%d"; sleep 1; date "+ %H:%M:%S"; sleep 1; done | fcd
```


#### `bash` - command line collaboration

```
# on raspi:
facade render term -mask=f

# on client:
facade exec -host raspi term -w 80 -h 25 bash
```


#### `frotz` - interactive fiction

```
# on raspi:
facade render term -shape slate -zoom .75

# on client:
facade exec -host raspi term -w 110 -h 30 frotz /path/to/hitchhikers_guide.z5
```


#### `man` - some manuals are quite pretty
```
# on raspi:
facade render lines -w 50 -shape crawl

# on client:
alias fcd='nc raspi 4045'
MANWIDTH=50 MANPAGER=cat man ssh \
 | while read line; do echo "$line"; sleep .9; done | fcd
```


#### `RFC` - protocol specifications in plaintext
```
# on raspi:
facade render lines -w 72 -shape rows

# on client:
alias fcd='nc raspi 4045'
curl -L https://tools.ietf.org/rfc/rfc792.txt \
 | while read -r line; do echo "$line"; sleep .9; done | fcd
```


#### `PHRACK` - your favourite hacking zine articles
```
# on raspi:
facade render lines -w 80 -shape roll

# on client:
alias fcd='nc raspi 4045'
curl -sL http://phrack.org/archives/tgz/phrack49.tar.gz \
 | tar xfz /dev/stdin --to-stdout ./14.txt \
 | while read -r line; do echo "$line"; sleep .9; done | fcd
```


#### `asciipr0n` - nudes older than the `<IMG>` tag
```
# on raspi:
facade render lines -w 80 -shape slate

# on client:
alias fcd='nc raspi 4045'
curl -sL https://www.asciipr0n.com/pr0n/pinups/pinup00.txt \
 | while read -r line; do echo "$line"; sleep .5; done | fcd
```


#### `.nfo` - demo scene release notez
```
# on raspi:
facade render lines -w=80 -vert=wave -mask=mask -font adore64

# on client:
alias fcd='nc raspi 4045'
curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
 | while read -r line; do echo "$line"; sleep .9; done | fcd
```


#### `parrot.live` - http streaming
```
# on raspi:
facade render term -shape slate

# on client:
facade exec term -w 50 -h 20 curl parrot.live
```

## Author

If you enjoy FACADE, tell us how you are using it at <facade@feedface.com>!

