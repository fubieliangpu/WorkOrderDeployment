package dci

type ConflictStatus int

const (
	VLANCONFLICT ConflictStatus = iota + 1
	VNICONFLICT
	VSINAMECONFLICT
	BRIGENUMCONFLICT
	//根据客户需求带宽量判断有没有足够多的空闲端口
	NOPORT
)
