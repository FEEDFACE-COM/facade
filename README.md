
# FACADE by FEEDFACE.COM

    
FACADE 



## Examples


### Setup Alias

    
    alias fcd='nc -N localhost 4045' # for linux
    alias fcd='nc localhost 4045'    # for mac/bsd







## Decorative / Wallpaper


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

    facade serve lines -vert drop -down -speed .2 -w 120 -h 8 -mask=mask
    tcpdump -i vlan5 -n -t -l -v dst port 53  | fcd

#### top

    facade exec term -w=64 -h=16 -vert=disk /usr/bin/top
    
    

## Collaborative / Interactive


# Manpages

    facade exec term -w=50 -h=20 man ssh




# frotz

    facade exec term -w=64 -h=16 /path/to/frotz /path/to/hitchhikers_guide.blb

    
    
    
# clear display

#    printf '\033[8;16;64t' # resize terminal


    clear | fcd
    

