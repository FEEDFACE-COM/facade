
package main

type FcdConfer   struct {}
func NewFcdConfer(host string, port uint, timeout float64) (*FcdConfer) { return new(FcdConfer) }
func (confer *FcdConfer) conf() { Info("confing... ") }

