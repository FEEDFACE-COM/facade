
# FACADE by FEEDFACE.COM

    
_FACADE_ is a creative coding tool that allows you to pipe any text on `/dev/stdout` directly onto the wall of your home / office / hackerspace. 

You will need:

* A Raspberry Pi running `facade serve`, and reachable by network
* A projector connected to the Raspberry Pi via HDMI
* An alias `alias fcd='nc -N raspberrypi 4045'` in your shell

Then just run `echo FOOBAR | fcd` on any machine that can reach the Raspberry Pi. This will render the word _FOOBAR_ on the projector for everyone to see. By default the text is displayed as a plain console grid, but FACADE includes more interesting styles in the form of vertex shaders that you can create and alter on the fly.

#### Motivation
the motivation is twofold:

1. nobody ever sees the stuff that happens in your shell, lets change that!

1. there's lots of really beautiful stuff available in 7-bit ascii on 80x25 chars



## Example Uses

### Informative / Status
You can use FACADE to show the status of your machines / networks.

* __top__ - system status

```
facade -q serve term -vert=wave -mask=mask &
facade exec -host raspberrypi term -h=24 -w=80 top -1   # run on server
```

* __tcpdump__ - live DNS queries

```
facade serve lines -w=80 -h=16 -buffer=32 -vert=roll -down
sudo tcpdump -i wlan0 -n -t -l dst port 53 | fcd        # run on router
```


* __access.log__ - live web requests

```
facade serve lines -w=120 -h=12 -buffer=16  -vert=vortex -mask=mask
tail -f /var/log/nginx/access.log | fcd                 # run on webserver
```

* __mtr__ - continuous trace route

```
facade -q serve term -vert=disk &
facade exec term -w=64 -h=16 sudo mtr -m 10 --displaymode 1 8.8.8.8
```

* __date__ - current date and time

```
facade serve -q lines -vert=wave -h=2 -w=10 -mask=mask -down -smooth=f -font=ocraext -zoom=.8 &
while true; do date "+%Y-%m-%d"; sleep 1; date "+ %H:%M:%S"; sleep 1; done | fcd
```



### Interactive / Collaborative

* __man__ - read manpages with your team

```
facade -q serve term -vert=def
facade exec term -w=50 -h=20 man ssh
```

* __bash__ - show your team what you are doing in a shell

```
facade -q serve term
facade exec term -w=80 -h=25 bash
```

* __frotz__ - play text adventures together

```
facade -q serve term -font=spacemono
facade exec term -w=64 -h=16 frotz /path/to/hitchhikers_guide.z5

```



### Decorative 



* __Phrack__ - your favourite hacking zine articles

```
facade -q serve lines -w=80 -h=25 -vert=roll &
curl -L http://phrack.org/archives/tgz/phrack49.tar.gz \
| tar xfz /dev/stdin ./14.txt --to-stdout \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
```

* __.nfo__ - demo scene release notes with 1337 ascii graphix

```
facade -q serve lines -w=80 -h=25 -vert=wave -font adore64 & 
curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
```

* __RFCs__ - internetworking specification in txt format

```
facade -q serve lines -w=72 -h=16 -vert=rows &
curl -L https://tools.ietf.org/rfc/rfc2460.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
```

* __man__ - some manpages are pretty too

```
facade -q serve lines -w=50 -h=20 -vert=crawl
MANWIDTH=50 MANPAGER=cat man ssh \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
```





#### Setup Alias

    
    alias fcd='nc -N localhost 4045' # for linux
    alias fcd='nc localhost 4045'    # for mac/bsd


## Development


To build _facade_ from source, clone the repository from github, then build with go. Specify a different output name to prevent filename clashes:

	git clone https://github.com/FEEDFACE-COM/facade.git
	cd facade
	go build -o facade-$(whoami)
	
To build a _facade_ binary with custom shaders or fonts, use GNU Make:

	git clone https://github.com/FEEDFACE-COM/facade.git
	cd facade
	make build


# Odds & Ends





    reset | fcd
    clear | fcd
    

# TODO
    daemonize
    unicode
    auth/tls
    
    
    
