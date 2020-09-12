# FACADE by FEEDFACE.COM
    
_FACADE_ is a creative coding tool that allows you to pipe any text on stdout directly onto the wall of your home / office / hackerspace. 

You will need:

- A Raspberry Pi running `facade serve`, and reachable by network
- A projector connected to the Raspberry Pi via HDMI
- An alias `alias fcd='nc -N raspberrypi 4045'` in your shell

Then just run `echo FOOBAR | fcd` on any machine that can reach the Raspberry Pi. This will render the text _FOOBAR_ onto the wall for everyone to see. While the text is displayed in plain console style by default, FACADE also supports custom styles in the form of OpenGLES shaders that you can create and alter on the fly.


--


The motivation for creating FACADE is twofold:

* We spend a lot of time in various consoles, terminals and shells, but most of the people around us never get to see the text we interact with.

* There is a lot of awesome plain text in the form of, RFCs, .nfo files, text adventures etc that looks even better when taken out of the confines of an xterm window and projected onto a wall.




## Compatibility 

FACADE server works on

- Raspbian Lite (no Desktop) on Raspberry Pi 2
- Raspbian Lite (no Desktop) on Raspberry Pi 3

FACADE client works on

- GNU/Linux on x86-64
- OpenBSD on x86-64
- Apple MacOS on x86-64




## Setup

1. Boot your Raspberry Pi into Raspbian, without the graphical X Server 

1. Download the latest release package for your platform(s):
	- Raspbian: `facade-x.y.z-linux-arm.tgz`
	- GNU/Linux: `facade-x.y.z-linux-amd64.tgz`
	- OpenBSD: `facade-x.y.z-openbsd-amd64.tgz`
	- Apple MacOS: `facade-x.y.z-darwin-amd64.tgz`

2. Extract the _facade_ binary from the release package:
	* `tar xfz facade-x.y.z-arch-os.tgz` 

3. On your Raspberry Pi, run FACADE server from console or ssh:
	* `./facade serve -d`
	
	You should now see the FACADE title screen on the HDMI display of the Raspberry Pi.

4. On your workstation, create the _fcd_ shell alias:  
   (replace _raspberrypi_ with the hostname or IP address of your Raspberry Pi)
	* Raspbian, Linux: `alias fcd='nc -N raspberrypi 4045'`
	* OpenBSD, MacOS: `alias fcd='nc raspberrypi 4045'`
	
5. Test whether you can send raw text to the FACADE server:
	* `whoami | fcd`

   You should now see your username on the HDMI display of the Raspberry Pi.

6. On your workstation, run FACADE client:  
   (replace _raspberrypi_ with the hostname or IP address of your Raspberry Pi)
	* `date --iso-8601 | ./facade pipe -host raspberrypi lines -font adore64`

	You should now see your username and the current date in a pixel font.

7. Explore FACADE options:

	* `./facade -h`
	* `./facade conf -h`
	* `./facade conf lines -h`




## Security

* FACADE provides no encryption, transport layer or otherwise.
* FACADE provides no authentication mechanism whatsoever.

The reasoning is that anyone able to reach the service and read the packets probably can see the output on the wall anyway, hence there is no focus on security at this time. **Please make sure to setup packet filters and network tunnels before sending any sensitive data to FACADE!**




## Custom Shaders

FACADE supports custom vertex and fragment shaders:

1. On your Raspberry Pi, create a _.facade/_ directory in your _$HOME_:
	* `mkdir -p ~/.facade/`
    
2. Download the default shader from the FACADE source repository:
	* `git archive --remote=https://github.com/FEEDFACE-COM/facade x.y.z \`
	`| tar x -C ~/.facade/ shader/grid/def.vert`

3. Copy the default shader to a new file _foobar.vert_:
	* `cp ~/.facade/shader/grid/def.vert ~/.facade/shader/grid/foobar.vert`

4. Instruct FACADE to use the new shader:
	* `./facade -d serve lines -vert foobar`

5. Edit the _~/.facade/shader/grid/foobar.vert_ file.
	* Add a line `pos.x*=sin(now); pos.y*=cos(now);` just before the line starting with `glPosition =` 	 

6. Save the _~/.facade/shader/grid/foobar.vert_ file.  
   You should see the effect of your changes in directly. 




## Examples

### Informative Use

FACADE can show you the live status of your machines, services and networks, eg:


* __top__ - system status

~~~
facade -q serve term -vert=wave -mask=mask &
facade exec -host raspberrypi term -h=24 -w=80 top -1   # run on server
~~~


* __tcpdump__ - live DNS queries

~~~
facade serve lines -w=80 -h=16 -buffer=32 -vert=roll -down
sudo tcpdump -i wlan0 -n -t -l dst port 53 | fcd        # run on router
~~~


* __access.log__ - live web requests

~~~
facade serve lines -w=120 -h=12 -buffer=16  -vert=vortex -mask=mask
tail -f /var/log/nginx/access.log | fcd                 # run on webserver
~~~


* __mtr__ - continuous trace route

~~~
facade -q serve term -vert=disk &
facade exec term -w=64 -h=16 sudo mtr -m 10 --displaymode 1 8.8.8.8
~~~


* __date__ - current date and time

~~~
facade serve -q lines -vert=wave -h=2 -w=10 -mask=mask -down -smooth=f -font=ocraext -zoom=.8 &
while true; do date "+%Y-%m-%d"; sleep 1; date "+ %H:%M:%S"; sleep 1; done | fcd
~~~



### Collaborative Use

You can use FACADE to look at text output together, ie one person directly interacts with a program while the other people in the room can observe and comment:


* __bash__ - show your team what exactly you are doing in your shell

~~~
facade -q serve term
facade exec -host raspberrypi term -w=80 -h=25 bash     # run on workstation
~~~


* __frotz__ - play text adventures as a group

~~~
facade -q serve term -font=spacemono
facade exec term -w=64 -h=16 frotz /path/to/hitchhikers_guide.z5
~~~



### Decorative Use

FACADE works very well if you just want to have some stylish text scrolling across your wall:


* __Phrack__ - your favourite hacking zine articles

~~~
facade -q serve lines -w=80 -h=25 -vert=roll &
curl -L http://phrack.org/archives/tgz/phrack49.tar.gz \
| tar xfz /dev/stdin ./14.txt --to-stdout \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


* __.nfo__ - demo scene release notes with 1337 ascii gfx

~~~
facade -q serve lines -w=80 -h=25 -vert=wave -font adore64 & 
curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


* __RFCs__ - internetworking specifications in plain text format

~~~
facade -q serve lines -w=72 -h=16 -vert=rows &
curl -L https://tools.ietf.org/rfc/rfc2460.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


* __man__ - some manpages are pretty too

~~~
facade -q serve lines -w=50 -h=20 -vert=crawl
MANWIDTH=50 MANPAGER=cat man git-rebase \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


--

If you enjoy FACADE, tell us at <facade@feedface.com> how you are using it!

