package go_grid

import "github.com/pyihe/util/math"

//地图在第一象限
type Map struct {
	width      int //地图宽度
	height     int //地图高度
	gridWidth  int //网格宽度
	gridHeight int //网格高度
	gridXCount int //横向grid的个数
	gridYCount int //纵向grid的个数

	grids []Grid //地图中的网格
	//TODO 灯塔
}

type Options func(g *Map)

func WithWidth(width int) Options {
	return func(m *Map) {
		m.width = width
	}
}

func WithHeight(height int) Options {
	return func(m *Map) {
		m.height = height
	}
}

func WithGridWidth(gWidth int) Options {
	return func(m *Map) {
		m.gridWidth = gWidth
	}
}

func WithGridHeight(gHeight int) Options {
	return func(m *Map) {
		m.gridHeight = gHeight
	}
}

func NewMap(options ...Options) *Map {
	var m = &Map{}
	for _, opt := range options {
		opt(m)
	}
	if m.width == 0 {
		m.width = defaultMapWidth
	}
	if m.height == 0 {
		m.height = defaultMapHeight
	}
	if m.gridWidth == 0 {
		m.gridWidth = defaultGridWidth
	}
	if m.gridHeight == 0 {
		m.gridHeight = defaultGridHeight
	}
	m.initGrid()
	return m
}

func (m *Map) Width() int {
	return m.width
}

func (m *Map) Height() int {
	return m.height
}

func (m *Map) GridWidth() int {
	return m.gridWidth
}

func (m *Map) GridHeight() int {
	return m.gridHeight
}

func (m *Map) GetGridById(id int) Grid {
	return m.grids[id]
}

func (m *Map) GetGridByGridXY(gridX, gridY int) Grid {
	gid := gridX + gridY*m.gridXCount
	if gid > 0 && gid < len(m.grids) {
		return m.grids[gid]
	}
	return nil
}

//这里的Point是Map中的实际坐标点
func (m *Map) GetGridByCoord(x, y int) Grid {
	if x < 0 || y < 0 {
		return nil
	}
	if x > m.width || y > m.height {
		return nil
	}
	gridX := x / m.gridWidth
	gridY := y / m.gridHeight
	return m.GetGridByGridXY(gridX, gridY)
}

//初始化网格
//网格
func (m *Map) initGrid() {
	var xCnt = m.width / m.gridWidth
	var yCnt = m.height / m.gridHeight
	if m.width%m.gridWidth != 0 {
		xCnt += 1
	}
	if m.height%m.gridHeight != 0 {
		yCnt += 1
	}

	m.gridXCount = xCnt
	m.gridYCount = yCnt
	m.grids = make([]Grid, xCnt*yCnt)

	for y := 0; y < yCnt; y++ {
		for x := 0; x < xCnt; x++ {
			gId := x + y*xCnt
			minX := x * m.gridWidth
			maxX := (x+1)*m.gridWidth - 1
			minY := y * m.gridHeight
			maxY := (y+1)*m.gridHeight - 1
			if maxX >= m.width {
				maxX = m.width - 1
			}
			if maxY >= m.height {
				maxY = m.height - 1
			}

			grid := &myGrid{
				id:       gId,
				width:    m.gridWidth,
				height:   m.gridHeight,
				x:        x,
				y:        y,
				minX:     minX,
				maxX:     maxX,
				minY:     minY,
				maxY:     maxY,
				entities: make(map[int64]Entity),
				points:   make([]Point, m.gridWidth*m.gridHeight),
			}
			m.grids[gId] = grid
		}
	}
}

//以focus为中心向外扩展
//level为遍历多少层网格, 为1时刚好为9宫格, 为2时则为16宫格
//fn表示待执行任务
func (m *Map) RangeEntity(focus Point, level int, fn func(entity Entity) error) {
	if level <= 0 {
		level = defaultLevelValue
	}
	focusGrid := m.GetGridByCoord(focus.X(), focus.Y())
	if focusGrid == nil {
		return
	}

	focusGrid.RangeEntity(fn)

	//计算以focus为中心, level为半径的遍历范围
	startX := math.MaxInt(focusGrid.GetGridX()-level, 0)
	endX := math.MinInt(focusGrid.GetGridX()+level, m.gridXCount-1)
	startY := math.MaxInt(focusGrid.GetGridY()-level, 0)
	endY := math.MaxInt(focusGrid.GetGridY()+level, m.gridYCount-1)

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			gId := x + y*m.gridXCount
			if gId == focusGrid.GetId() {
				continue
			}
			grid := m.GetGridById(gId)
			if grid == nil {
				continue
			}
			grid.RangeEntity(fn)
		}
	}
}

//计算地图中两个实际坐标之间的网格距离
func (m *Map) GetGridDistance(src, target Point) []Grid {
	var result []Grid

	srcGrid := m.GetGridByCoord(src.X(), src.Y())
	targetGrid := m.GetGridByCoord(target.X(), target.Y())
	if srcGrid == nil || targetGrid == nil {
		return result
	}

	sx, sy, tx, ty := srcGrid.GetGridX(), srcGrid.GetGridY(), targetGrid.GetGridX(), targetGrid.GetGridY()
	xChange, yChange := tx-sx, ty-sy

	var xRange = 1
	var yRange = 1

	if xChange > 0 { // 右移
		for x := 0; x < xChange; x++ {
			for y := sy - yRange; y <= sy+yRange; y++ {
				result = append(result, m.GetGridByCoord(x, y)) // 左
			}
			for y := ty - yRange; y <= ty+yRange; y++ {
				result = append(result, m.GetGridByCoord(x, y)) // 右
			}
		}
	} else { // 左移
		for x := 0; x > xChange; x-- {
			for y := ty - yRange; y <= ty+yRange; y++ {
				result = append(result, m.GetGridByCoord(x, y)) // 左
			}
			for y := sy - yRange; y <= sy+yRange; y++ {
				result = append(result, m.GetGridByCoord(x, y)) // 右
			}
		}
	}
	if yChange < 0 { // 上移
		for y := 0; y > yChange; y-- {
			for x := tx - xRange; x <= tx+xRange; x++ {
				result = append(result, m.GetGridByCoord(x, y)) // 上
			}
			for x := sx - xRange; x <= sx+xRange; x++ {
				result = append(result, m.GetGridByCoord(x, y)) // 下
			}
		}
	} else { // 下移
		for y := 0; y < yChange; y++ {
			for x := sx - xRange; x <= sx+xRange; x++ {
				result = append(result, m.GetGridByCoord(x, y)) // 上
			}
			for x := tx - xRange; x <= tx+xRange; x++ {
				result = append(result, m.GetGridByCoord(x, y)) // 下
			}
		}
	}
	return result
}
