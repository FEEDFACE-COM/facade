
package main

import (
    "fmt"
    conf "./conf"
    log "./log"
    "golang.org/x/image/math/fixed"
)


func mask(in fixed.Int26_6) int { return (0xffffffff & int(in)) }
func str(s string, in fixed.Int26_6) string { return fmt.Sprintf("%s\t0x%08x\t%v\t %d ≤ %d ≤ %d",s,mask(in),in.String(),in.Floor(),in.Round(),in.Ceil()) }

func (tester *Tester) testFixed(config *conf.Config) error {
    log.Info("testing fixed...")    
    
    one_one_fourth := fixed.Int26_6( 1<<6 + 1<<4 )
    one := fixed.Int26_6(1 << 6)
    two := fixed.Int26_6(2 << 6)
    three := fixed.Int26_6(3 << 6)
    four := fixed.Int26_6(4 << 6 )

    half := fixed.Int26_6( 0x20 )
    fourth := fixed.Int26_6( 0x10 )

    four_halves := half.Mul(four)
    four_fourths := fourth.Mul(four)
    fourth_fours := four.Mul(fourth)
    five := one_one_fourth.Mul(four)
    
    log.Info(str("¼",fourth))
    log.Info(str("½",half))
    log.Info(str("1",one))
    log.Info(str("¼x4",four_fourths))
    log.Info(str("4x¼",fourth_fours))
    log.Info(str("1¼",one_one_fourth))
    log.Info(str("2",two))
    log.Info(str("4x½",four_halves))
    log.Info(str("3",three))
    log.Info(str("4",four))
    log.Info(str("4x1¼",five))
    
    return nil
}
