package go_grid

import "sync"

//点
type Point interface {
	X() int //横坐标
	Y() int //纵坐标
}

//实体
type Entity interface {
	GetId() int64 //获取实体ID
}

//网格
type Grid interface {
	//grid属性相关
	GetId() int
	GetGridX() int
	GetGridY() int

	//grid包含的实体相关
	AddEntity(entity Entity) error
	GetEntity(id int64) (Entity, error)
	RemoveEntity(id int64)
	RangeEntity(fn func(entity Entity) error)

	//grid中的坐标相关
	IsInGrid(realX, realY int) bool
	SetPoint(p Point)
	GetPoint(realX, realY int) Point
}

//网格
type myGrid struct {
	sync.Mutex                  //锁
	id         int              //id
	width      int              //网格宽度
	height     int              //网格高度
	x          int              //以网格为最小单位时, 网格的横坐标
	y          int              //以网格为最小单位时, 网格的纵坐标
	minX       int              //网格最小横坐标
	maxX       int              //网格最大横坐标
	minY       int              //网格最小纵坐标
	maxY       int              //网格最大纵坐标
	entities   map[int64]Entity //网格内的实体, 比如: NPC, Monster, Item...
	points     []Point          //网格内的坐标点
}

func (g *myGrid) GetId() int {
	return g.id
}

func (g *myGrid) GetGridX() int {
	return g.x
}

func (g *myGrid) GetGridY() int {
	return g.y
}

func (g *myGrid) AddEntity(entity Entity) error {
	g.Lock()
	defer g.Unlock()

	if _, ok := g.entities[entity.GetId()]; ok {
		return ErrAlreadyExistEntity
	}
	g.entities[entity.GetId()] = entity
	return nil
}

func (g *myGrid) GetEntity(id int64) (Entity, error) {
	if g.entities == nil {
		return nil, ErrNoEntity
	}
	if e, ok := g.entities[id]; ok {
		return e, nil
	}
	return nil, ErrNoEntity
}

func (g *myGrid) RemoveEntity(id int64) {
	g.Lock()
	defer g.Unlock()

	delete(g.entities, id)
}

//判断x, y对应的点是否在网格内
func (g *myGrid) IsInGrid(realX, realY int) bool {
	x := realX - g.minX
	y := realY - g.minY
	return x >= 0 && x < g.width && y >= 0 && y < g.height
}

func (g *myGrid) SetPoint(p Point) {
	if !g.IsInGrid(p.X(), p.Y()) {
		return
	}
	x := p.X() - g.minX
	y := p.Y() - g.minY
	g.points[x+y*g.width] = p
}

func (g *myGrid) GetPoint(realX, realY int) Point {
	if !g.IsInGrid(realX, realY) {
		return nil
	}
	x := realX - g.minX
	y := realY - g.minY
	return g.points[x+y*g.width]
}

//遍历grid中的实体, 并执行entity
func (g *myGrid) RangeEntity(fn func(entity Entity) error) {
	var err error
	for _, en := range g.entities {
		if err = fn(en); err != nil {
			return
		}
	}
}
