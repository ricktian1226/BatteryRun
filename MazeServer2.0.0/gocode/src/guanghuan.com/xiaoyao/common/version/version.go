package xyversion

import (
	"fmt"
)

type Version int32

func (ver *Version) String() (str string) {
	v := int32(*ver)
	z := v % 100
	y := (v / 100) % 100
	x := v / 10000
	str = fmt.Sprintf("%d.%02d.%02d", x, y, z)

	return
}

func New(x int32, y int32, z int32) (ver *Version) {
	ver = new(Version)
	*ver = Version(x*10000 + y*100 + z)

	return
}
func (ver *Version) Set(x, y, z int32) {
	*ver = Version(x*10000 + y*100 + z)
}
func (ver *Version) LowerThan(ver_right Version) bool {
	return int32(*ver) < int32(ver_right)
}

func (ver *Version) LargerThan(ver_right Version) bool {
	return int32(*ver) > int32(ver_right)
}

func (ver *Version) EqualTo(ver_right Version) bool {
	return int32(*ver) == int32(ver_right)
}
