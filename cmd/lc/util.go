package main

import "github.com/lab5e/lospan/pkg/pb/lospan"

func newPtr[T int | uint32 | int32 | bool | float32 | string | lospan.DeviceState](v T) *T {
	ret := new(T)
	*ret = v
	return ret
}

func ellipsisString(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
