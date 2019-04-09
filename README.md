
fresh install:

[folkert@korova ~]$ ./bin/fcd-linux-arm -d 
./bin/fcd-linux-arm: error while loading shared libraries: libGL.so.1: cannot open shared object file: No such file or directory



=>


sudo apt-get install mesa-utils




=>

delta conf[]

   _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _        
  |_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|  
  |   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |  

init renderer[~/src/gfx/facade] conf[]
* failed to add service - already in use?





=>

gfx/axis.go:7:5: cannot find package "github.com/go-gl/mathgl/mgl32" in any of:
	/usr/lib/go-1.7/src/github.com/go-gl/mathgl/mgl32 (from $GOROOT)
	/home/folkert/go/src/github.com/go-gl/mathgl/mgl32 (from $GOPATH)
gfx/font.go:14:5: cannot find package "github.com/golang/freetype" in any of:
	/usr/lib/go-1.7/src/github.com/golang/freetype (from $GOROOT)
	/home/folkert/go/src/github.com/golang/freetype (from $GOPATH)
gfx/font.go:15:5: cannot find package "github.com/golang/freetype/truetype" in any of:
	/usr/lib/go-1.7/src/github.com/golang/freetype/truetype (from $GOROOT)
	/home/folkert/go/src/github.com/golang/freetype/truetype (from $GOPATH)
gfx/font.go:13:5: cannot find package "golang.org/x/image/font" in any of:
	/usr/lib/go-1.7/src/golang.org/x/image/font (from $GOROOT)
	/home/folkert/go/src/golang.org/x/image/font (from $GOPATH)
render_linux_arm.go:15:5: cannot find package "src.feedface.com/gfx/piglet" in any of:
	/usr/lib/go-1.7/src/src.feedface.com/gfx/piglet (from $GOROOT)
	/home/folkert/go/src/src.feedface.com/gfx/piglet (from $GOPATH)
gfx/axis.go:9:5: cannot find package "src.feedface.com/gfx/piglet/gles2" in any of:
	/usr/lib/go-1.7/src/src.feedface.com/gfx/piglet/gles2 (from $GOROOT)
	/home/folkert/go/src/src.feedface.com/gfx/piglet/gles2 (from $GOPATH)
Makefile:76: recipe for target 'fcd-linux-arm' failed


=>


go get -v

