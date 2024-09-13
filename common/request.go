package common

type UPDATE_MODE int

const (
	//全量更新
	UPDATE_MODE_PUT = iota
	//部分更新
	UPDATE_MODE_PATCH
)
