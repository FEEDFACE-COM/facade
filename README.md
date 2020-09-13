# FACADE by FEEDFACE.COM
    
_FACADE_ is a creative coding tool that allows you to pipe any text on stdout directly onto the wall of your home / office / hackerspace. 

You will need:

- A Raspberry Pi running `facade serve`, and reachable by network
- A projector connected to the Raspberry Pi via HDMI
- An alias `alias fcd='nc -N raspberrypi 4045'` in your shell

Then just run `echo FOOBAR | fcd` on any machine that can reach the Raspberry Pi. This will render the text _FOOBAR_ onto the wall for everyone to see. The text is displayed in plain console style by default, FACADE also supports custom styles in the form of OpenGLES shaders that you can create and alter on the fly.


----


The motivation for creating FACADE is twofold:

* We spend a lot of time in various consoles, terminals and shells, but most of the people around us never get to see the text we interact with.

* There is a lot of awesome plain text in the form of RFCs, .nfo files, text adventures etc that looks even better when taken out of the confines of an xterm window and projected onto a wall.




## Requirements 

FACADE server works on

- Raspbian Lite on Raspberry Pi 2
- Raspbian Lite on Raspberry Pi 3

FACADE client works on

- GNU/Linux on x86-64
- OpenBSD on x86-64
- Apple MacOS on x86-64




## Setup

### Setup FACADE Server on Raspberry Pi

1. Prepare your Raspberry Pi configuration:
    - Use `raspi-config` to set the memory available to the GPU to _256 MB_
    - Use `raspi-config` to select the _Legacy non-GL desktop driver_
    - Make sure the X Window System is not running

2. Download the latest release package:
	- Raspbian Lite: `facade-x.y.z-linux-arm.tgz`

3. Extract the _facade_ binary from the release package:
	- `tar xfz facade-x.y.z-linux-arm.tgz` 

4. Run FACADE server from a console or ssh shell:
	- `./facade serve -d`
	
	You should now see the FACADE title screen on the HDMI display of the Raspberry Pi.

5. On your workstation, create the _fcd_ shell alias:  
   (replace _raspberrypi_ with the hostname or IP address of your Raspberry Pi)
	- Raspbian, Linux: `alias fcd='nc -N raspberrypi 4045'`
	- OpenBSD, MacOS: `alias fcd='nc raspberrypi 4045'`
	
6. Test whether you can send raw text to the FACADE server:
	- `whoami | fcd`

   You should now see your username on the HDMI display of the Raspberry Pi.


### Setup FACADE Client on your workstation

1. Download the latest release package for your platform:
	- GNU/Linux: `facade-x.y.z-linux-amd64.tgz`
	- OpenBSD: `facade-x.y.z-openbsd-amd64.tgz`
	- Apple MacOS: `facade-x.y.z-darwin-amd64.tgz`

2. Extract the _facade_ binary from the release package:
	- `tar xfz facade-x.y.z-os-arch.tgz` 

3. Run FACADE client:  
   (replace _raspberrypi_ with the hostname or IP address of your Raspberry Pi)
	- `date --iso-8601 | ./facade pipe -host raspberrypi lines -font adore64`

	You should now see your username and the current date in a pixel font.

4. Explore FACADE options:

	- `./facade -h`
	- `./facade conf -h`
	- `./facade conf lines -h`




## Security

- FACADE provides no encryption, transport layer or otherwise.
- FACADE provides no authentication mechanism whatsoever.

The reasoning is that anyone able to reach the service and read the packets probably can see the output on the wall anyway, hence there is no focus on security at this time. **Please make sure to setup packet filters and network tunnels before sending any sensitive data to FACADE!**




## Custom Shaders

FACADE supports custom vertex and fragment shaders:

1. On your Raspberry Pi, create a _.facade/_ directory in your _$HOME_:
	- `mkdir -p ~/.facade/`
    
2. Download the default shader from the FACADE source repository:
	- `git archive --remote=https://github.com/FEEDFACE-COM/facade HEAD | tar x -C ~/.facade/ shader/grid/def.vert`

3. Copy the default shader to a new file _foobar.vert_:
	- `cp ~/.facade/shader/grid/def.vert ~/.facade/shader/grid/foobar.vert`

4. Instruct FACADE to use the new shader:
	- `./facade -d serve lines -vert foobar`

5. Edit the _~/.facade/shader/grid/foobar.vert_ file.
	- Try adding `pos.x*=sin(now); pos.y*=cos(now);` just before the line starting with `glPosition` 	 

6. Save the _~/.facade/shader/grid/foobar.vert_ file.  
   - You should see the effect of your changes. 

See <https://www.khronos.org/opengles/sdk/docs/reference_cards/OpenGL-ES-2_0-Reference-card.pdf> for shader syntax and available functions.



## Examples

### Informative Use

FACADE can show you the live status of your machines, services and networks, eg:


#### `top` - system status
~~~
facade -q serve term -vert=wave -mask=mask &
facade exec -host raspberrypi term -h=24 -w=80 top -1   # run on server
~~~


#### `tcpdump` - live DNS queries
~~~
facade serve lines -w=80 -h=16 -buffer=32 -vert=roll -down
sudo tcpdump -i wlan0 -n -t -l dst port 53 | fcd        # run on router
~~~


#### `access.log` - live web requests
~~~
facade serve lines -w=120 -h=12 -buffer=16  -vert=vortex -mask=mask
tail -f /var/log/nginx/access.log | fcd                 # run on webserver
~~~


#### `mtr` - continuous trace route
~~~
facade -q serve term -vert=disk &
facade exec term -w=64 -h=16 sudo mtr -m 10 --displaymode 1 8.8.8.8
~~~


#### `date` - current date and time
~~~
facade serve -q lines -vert=wave -h=2 -w=10 -mask=mask -down -smooth=f -font=ocraext -zoom=.8 &
while true; do date "+%Y-%m-%d"; sleep 1; date "+ %H:%M:%S"; sleep 1; done | fcd
~~~



### Collaborative Use

You can use FACADE to look at text output together, ie one person directly interacts with a program while the other people in the room can observe and comment:


#### `bash` - show your team what exactly you are doing in your shell
~~~
facade -q serve term
facade exec -host raspberrypi term -w=80 -h=25 bash     # run on workstation
~~~


#### `frotz` - play text adventures as a group
~~~
facade -q serve term -font=spacemono
facade exec term -w=64 -h=16 frotz /path/to/hitchhikers_guide.z5
~~~



### Decorative Use

FACADE works very well if you just want to have some stylish text scrolling across your wall:


#### `man` - some manpages are quite pretty :)
~~~
facade -q serve lines -w=50 -h=20 -vert=crawl
MANWIDTH=50 MANPAGER=cat man git-rebase \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


#### `RFCs` - internetworking specifications in plain text format
~~~
facade -q serve lines -w=72 -h=16 -vert=rows &
curl -L https://tools.ietf.org/rfc/rfc2460.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


#### `PHRACK` - your favourite hacking zine articles
~~~
facade -q serve lines -w=80 -h=25 -vert=roll &
curl -L http://phrack.org/archives/tgz/phrack49.tar.gz \
| tar xfz /dev/stdin ./14.txt --to-stdout \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


#### `.nfo` - demo scene release notes with 1337 ascii art
~~~
facade -q serve lines -w=80 -h=25 -vert=wave -mask=mask -font adore64 & 
curl -L https://content.pouet.net/files/nfos/00012/00012031.txt \
| while read -r line; do echo "$line"; sleep .9; done \
| fcd
~~~


----

If you enjoy FACADE, tell us how you are using it at <facade@feedface.com>!

