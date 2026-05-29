package game

import (
	"math"
)

type SpatialGrid[T any] struct {
	cells    map[[2]int][]T
	cellSize float64
}

func NewSpatialGrid[T any](cellSize float64) *SpatialGrid[T] {
	return &SpatialGrid[T]{
		cells:    make(map[[2]int][]T),
		cellSize: cellSize,
	}
}

func (sg *SpatialGrid[T]) Clear() {
	clear(sg.cells)
}

func (sg *SpatialGrid[T]) Insert(pos Vec2, item T) {
	k := sg.key(pos)
	sg.cells[k] = append(sg.cells[k], item)
}

func (sg *SpatialGrid[T]) GetNeighbors(pos Vec2, radius float64) []T {
	var result []T
	k := sg.key(pos)
	spread := int(math.Ceil(radius / sg.cellSize))
	for dx := -spread; dx <= spread; dx++ {
		for dy := -spread; dy <= spread; dy++ {
			key := [2]int{k[0] + dx, k[1] + dy}
			result = append(result, sg.cells[key]...)
		}
	}

	return result
}

func (sg *SpatialGrid[T]) GetInRect(minX, minY, maxX, maxY float64) []T {
	var result []T
	startX := int(minX / sg.cellSize)
	startY := int(minY / sg.cellSize)
	endX := int(maxX / sg.cellSize)
	endY := int(maxY / sg.cellSize)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {
			result = append(result, sg.cells[[2]int{x, y}]...)
		}
	}

	return result
}

func (sg *SpatialGrid[T]) key(pos Vec2) [2]int {
	return [2]int{
		int(pos.X / sg.cellSize),
		int(pos.Y / sg.cellSize),
	}
}
