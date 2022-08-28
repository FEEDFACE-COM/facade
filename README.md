# FACADE by FEEDFACE.COM
    

### Informative Use

FACADE can show you the live status of your machines, services and networks, eg:


#### system status: `top`
```

facade serve term -shape wave &
facade exec -host raspi term -w 80 -h 25 top -1


```


#### network traffic: `tcpdump`
```

facade serve words -shape field -n 32 -life 4 -mark 1 -shuffle &
sudo tcpdump -i wlan0 -n -l -t tcp[tcpflags]=tcp-syn \
 | egrep --line-buffered -o "> [0-9]+\.[0-9]+.[0-9]+\.[0-9]+"  \
 | awk "//{print \$2;fflush}"  \
 | nc -D raspi 4045                                  


```


#### access logs: `tail`
```

facade serve lines -shape disk  -w 150 -h 12 &
tail -f /var/log/nginx/access.log \
 | nc -D raspi 4045


```

#### trace route: `mtr`
```

facade serve term -shape vortex &
facade exec term -w 120 -h 16 sudo mtr -m 10 --displaymode 1 wikipedia.org


```


#### clock: `date`
```

facade serve chars -shape moebius -w 64 -speed .5 -font spacemono &
while true; do date +"%Y-%m-%dT%H:%M:%S%z"; sleep 1; done \
 | nc -D raspi 4056


```

```

facade -d serve lines -shape wave -h 2 -w 10 -down -font ocraext -zoom .8 -smooth=f &
while true; do \
  date "+%Y-%m-%d"; sleep 1; date "+ %H:%M:%S"; sleep 1; \
done \
 | nc -D raspi 4056

```


	


### Collaborative Use

You can use FACADE to look at text output together, ie one person directly interacts with a program while the other people in the room can observe and comment:


#### `bash` - show your team what exactly you are doing in your shell

```

facade serve term &
facade exec -host raspi term -w 80 -h 25 bash


```


#### `frotz` - play text adventures on your wall

```
facade serve -dir=. term -shape slate &
facade exec term -w 110 -h 30 frotz /path/to/hitchhikers_guide.z5
```




### Decorative Use

FACADE works very well if you just want to have some stylish text scrolling across your wall:


#### `man` - some manpages are quite pretty :)
```
facade -q serve lines -w=50 -vert=crawl &
MANWIDTH=50 MANPAGER=cat man ssh \
| while read -r line; do echo "$line"; sleep .9; done | fcd
```


#### `rfc` - internetworking specifications in plain text format
```
facade serve -dir=. lines -w=72 -shape rows
curl -L https://tools.ietf.org/rfc/rfc792.txt \
| while read -r line; do echo "$line"; sleep .9; done | fcd
```


#### `PHRACK` - your favourite hacking zine articles
```
facade serve -dir=. lines -w=80 -shape roll &
curl -sL http://phrack.org/archives/tgz/phrack49.tar.gz \
| tar xfz /dev/stdin --to-stdout ./14.txt \
| while read -r line; do echo "$line"; sleep .9; done | fcd
```


#### `pr0n` - online nudes before the `<IMG>` tag
```
facade serve -dir=. lines -w 80 -shape slate 
curl -sL https://www.asciipr0n.com/pr0n/pinups/pinup00.txt \
| while read -r line; do echo "$line"; sleep .9; done | fcd

```


#### `.nfo` - demo scene release notes with 1337 ascii art
```
facade -q serve lines -w=80 -vert=wave -mask=mask -font adore64 &
curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
| while read -r line; do echo "$line"; sleep .9; done | fcdx
```


#### curl parrot.live
```
facade serve term 
facade exec term -w 40 -h 20 curl parrot.live

```

----


If you enjoy FACADE, tell us how you are using it at <facade@feedface.com>!

