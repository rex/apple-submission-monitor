package tui

const (
	headerHeight = 3
	footerHeight = 2
)

type rect struct {
	x      int
	y      int
	width  int
	height int
}

func (r rect) contains(x int, y int) bool {
	return x >= r.x && x < r.x+r.width && y >= r.y && y < r.y+r.height
}

type gridLayout struct {
	width   int
	height  int
	columns int
	rows    int
	count   int
	rects   []rect
}

func calculateGrid(width int, height int, count int) gridLayout {
	availableHeight := max(1, height-headerHeight-footerHeight)
	if count <= 0 || width <= 0 {
		return gridLayout{width: width, height: availableHeight}
	}

	columns := chooseColumns(width, availableHeight, count)
	rows := (count + columns - 1) / columns
	layout := gridLayout{
		width:   width,
		height:  availableHeight,
		columns: columns,
		rows:    rows,
		count:   count,
		rects:   make([]rect, 0, count),
	}
	for index := range count {
		layout.rects = append(layout.rects, layout.cell(index))
	}
	return layout
}

func (g gridLayout) cell(index int) rect {
	column := 0
	row := index
	left := column * g.width / g.columns
	right := (column + 1) * g.width / g.columns
	top := headerHeight + row*g.height/g.rows
	bottom := headerHeight + (row+1)*g.height/g.rows
	return rect{x: left, y: top, width: right - left, height: bottom - top}
}

func (g gridLayout) hit(x int, y int) int {
	for index, area := range g.rects {
		if area.contains(x, y) {
			return index
		}
	}
	return -1
}

func chooseColumns(_, _, _ int) int {
	return 1
}
