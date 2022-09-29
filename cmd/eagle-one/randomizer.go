package main

import "math/rand"

// Randomizer triggers a function a percentage of the time
type Randomizer struct {
	percent int
}

// NewRandomizer creates a new randomizer
func NewRandomizer(percentage int) *Randomizer {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	return &Randomizer{percentage}
}

// Maybe runs something a percentage of the time
func (t *Randomizer) Maybe(something func()) {
	if t.Now() {
		something()
	}
}

// Now returns true if the condition should trigger
func (t *Randomizer) Now() bool {
	return rand.Intn(100) < t.percent
}
