
package gfx
import ( "math" )



var PI float64  = 3.1415926535897932384626433832795028841971693993751058209749445920
var TAU float64 = 6.2831853071795864769252867665590057683943387987502116419498891840





func rads(deg float64) float64 { return TAU * deg/360. }
func degs(rad float64) float64 { return 360. * rad/TAU }

func ease0(x,f,p float64) float64 { return       math.Cos( f*x + PI/2. + p ) }
func ease1(x,f,p float64) float64 { return 0.5 * math.Cos( f*x + PI/2. + p ) + 0.5 } 

func max(a,b float64) float64 { if a>=b { return a }; return b }
func min(a,b float64) float64 { if a<=b { return a }; return b }
func abs(a float64) float64 { if a < 0 { return -a }; return a }
