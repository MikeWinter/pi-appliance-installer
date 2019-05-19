// +build !arm

package dev

const (
	_ Mode = iota
	BLOCK
	RW_USR
	_
	RW_GRP
)
