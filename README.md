
# FACADE by FEEDFACE.COM

    
_FACADE_ is a creative coding tool that allows you to pipe any text on `/dev/stdout` directly onto the wall of your home / office / hackerspace. 

You will need:

* A Raspberry Pi running `facade serve`, and reachable by network
* A projector connected to the Raspberry Pi via HDMI
* An alias `alias fcd='nc -N raspberrypi 4045'` in your shell

Then just run `echo FOOBAR | fcd` on any machine that can reach the Raspberry Pi. This will render the word _FOOBAR_ on the projector for everyone to see. By default the text is displayed as a plain console grid, but _FACADE_ includes more interesting styles in the form of vertex shaders that you can create and alter on the fly.




## Example Uses

### Informative Uses
### Collaborative Uses

### Decorative Use



* **Phrack** - your favorite ezine articles

```
facade -q serve lines -w=80 -h=25 -vert=roll &
curl -L http://phrack.org/archives/tgz/phrack49.tar.gz \
| tar xfz /dev/stdin ./14.txt --to-stdout \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
```

* **.nfo** - 1337 demo scene release notes

```
facade -q serve lines -w=80 -h=25 -vert=wave -font adore64 & 
curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
```

## Setup

### Alias `fcd`




### Setup


* A **Raspberry Pi** running `facade serve`
* Ideally, a **Projector** plugged into the HDMI port of the Raspberry Pi






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


## Example Uses









### Decorative / Wallpaper


#### Phrack


    facade conf lines -w=80 -h=25 -vert=roll
    curl -L http://phrack.org/archives/tgz/phrack49.tar.gz \
     | tar xfz /dev/stdin ./14.txt -O \
     | while read -r line; do echo "$line"; sleep .9; done \
     | fcd


#### RFC

    facade conf lines -w=72 -h=16 -vert=rows
    curl -L https://tools.ietf.org/rfc/rfc2460.txt \
     | while read -r line; do echo "$line"; sleep .9; done \
     | fcd

#### .nfo 

    facade conf lines -w=80 -h=25 -vert=wave -font adore64
    curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
     | while read -r line; do echo "$line"; sleep .9; done \
     | fcd
    
# Manpages    
    
    facade conf lines -w=50 -h=20 -vert=crawl
    export MANWIDTH=50 MANPAGER=cat
    man ssh | while read -r line; do echo "$line"; sleep .9; done | fcd





## Informative / Status




#### Access Logs

    facade serve lines -w=120 -h=8 -vert=disk -mask=mask &
    tail -f /var/log/nginx/access.log | fcd


#### Clock

    facade serve lines -vert=wave -h=2 -w=10 -mask=mask -down -smooth=f -font ocraext -zoom .8
    while true; do date "+%Y-%m-%d"; sleep 1; date "+ %H:%M:%S"; sleep 1; done | fcd

    
#### mtr
    
    facade serve term -vert=disk
    facade exec term  -w=64 -h=16 sudo mtr -m 10 --displaymode 2 example.com


#### tcpdump

    facade serve lines -vert vortex -down -speed .2 -w 120 -h 8 -mask=mask
    tcpdump -i vlan5 -n -t -l -v dst port 53  | fcd

#### top

    facade exec term -w=64 -h=16 -vert=disk /usr/bin/top
    
    

## Collaborative / Interactive


# Manpages

    facade exec term -w=50 -h=20 man ssh




# frotz

    facade exec term -w=64 -h=16 /path/to/frotz /path/to/hitchhikers_guide.blb

    
    
    
    reset | fcd
    clear | fcd
    

# TODO
    daemonize
    unicode
    auth/tls
    
    
    
