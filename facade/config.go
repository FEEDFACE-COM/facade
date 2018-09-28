
package facade


var DEFAULT_MODE Mode = TEST


type Mode string
const (
    GRID  Mode = "grid"
    LINES Mode = "lines"
    WORD  Mode = "word"
    CHAR  Mode = "char"   
    TEST  Mode = "test" 
)

var Modes = []Mode{GRID,LINES,TEST}




var DIRECTORY = "/home/folkert/src/gfx/facade/asset/"

