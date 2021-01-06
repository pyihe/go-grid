package go_grid

import (
	"github.com/pyihe/util/errors"
)

const (
	defaultMapWidth   = 100 //地图默认宽度
	defaultMapHeight  = 100 //地图默认高度
	defaultGridWidth  = 10  //网格默认宽度
	defaultGridHeight = 10  //网格默认高度
	defaultLevelValue = 1   //网格默认遍历层数
)

var (
	ErrAlreadyExistEntity = errors.NewError(-1, "entity already exist")
	ErrNoEntity           = errors.NewError(-1, "not found entity")
)
