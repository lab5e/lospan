package main

func newPtr[T int | bool | float32 | string](v T) *T {
	ret := new(T)
	*ret = v
	return ret
}
