
package math32

import ( "math" )



var PI float32  = 3.1415926535897932384626433832795028841971693993751058209749445920
var TAU float32 = 6.2831853071795864769252867665590057683943387987502116419498891840



func Rads(deg float32) float32 { return TAU * deg/360. }
func Degs(rad float32) float32 { return 360. * rad/TAU }



func Max(a,b float32) float32 { if a>=b { return a }; return b }
func Min(a,b float32) float32 { if a<=b { return a }; return b }
func Abs(a float32) float32 { if a < 0 { return -a }; return a }

func Cos(x float32) float32 { return float32( math.Cos( float64(x) ) ) }
func Floor(x float32) float32 { return float32( math.Floor( float64(x) ) ) }

func Clamp(x float32) float32 { if x < 0.0 { return 0.0 } else if x > 1.0 { return 1.0 } else { return x }   }


func EaseInEaseOut(x float32) float32 {
    return  -0.5 * Cos( x * PI  ) + 0.5
}