package main

import "math/rand"

// TheRandomizer triggers a function a percentage of the time
type TheRandomizer struct {
	percent int
}

// NewRandomizer creates a new randomizer
func NewRandomizer(percentage int) *TheRandomizer {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	return &TheRandomizer{percentage}
}

// Maybe runs something a percentage of the time
func (t *TheRandomizer) Maybe(something func()) {
	if t.Now() {
		something()
	}
}

// Now returns true if the condition should trigger
func (t *TheRandomizer) Now() bool {
	return rand.Intn(100) < t.percent
}
