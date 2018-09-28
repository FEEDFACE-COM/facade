
// +build !linux !arm

package main

import ( facade "./facade" )

type Tester struct {}
func NewTester(_ string) *Tester { return &Tester{} }
func (tester *Tester) Test(str string) { }
func (tester *Tester) Configure(config *facade.Config) {}

