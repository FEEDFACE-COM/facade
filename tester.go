
// +build !linux !arm

package main

import (
    conf "./conf"
)

type Tester struct {}
func NewTester() *Tester { return &Tester{} }
func (tester *Tester) Test(str string) { }
func (tester *Tester) Configure(config *conf.Config) {}

