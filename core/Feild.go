package core

import (
	"errors"
	"fmt"
)

type Field struct {
	matrix [][]FieldObject
	size   int
}

func NewField(size int) *Field {
	result := new(Field)
	result.size = size
	result.matrix = make([][]FieldObject, size)
	for i := range result.matrix {
		result.matrix[i] = make([]FieldObject, size)
	}

	filler := func() FieldObject {
		return NewEmptyPlace()
	}
	FillFieldWith(result.matrix, filler)
	return result
}

func (f *Field) AddObject(p FieldPoint, obj FieldObject) bool {
	if !f.checkCorectPoint(p) {
		return false
	}
	if !IsReplaceble(f.matrix[p.x][p.y]) {
		return false
	}
	f.matrix[p.x][p.y] = obj
	return true
}

func (f *Field) LookAt(p FieldPoint) (FieldObject, error) {
	if !f.checkCorectPoint(p) {
		return nil, errors.New("invalid location")
	}
	return f.matrix[p.x][p.y], nil
}

func (f *Field) GetAllWithType(t ObjType) []*ObjLocator {

	return f.GetAllWithTypeInSquare(GenFuncForObjType(t), FieldPoint{0, 0}, f.size)
}

func (f *Field) checkCorectPoint(p FieldPoint) bool {
	if p.x < 0 || p.y < 0 || p.x > f.size || p.y > f.size {
		return false
	}
	return true
}

func (f *Field) GetAllWithTypeInSquare(comp func(FieldObject) bool, center FieldPoint, size int) []*ObjLocator {
	result := make([]*ObjLocator, 0)

	topLeft := NewFieldPoint(center.x-size/2, center.y-size/2)
	min := func(x int, y int) int {
		if x > y {
			return y
		}
		return x
	}
	max := func(x int, y int) int {
		if x > y {
			return x
		}
		return y
	}
	fromX := max(0, topLeft.x)
	fromY := max(0, topLeft.y)
	toX := min(f.size, topLeft.x+size)
	toY := min(f.size, topLeft.y+size)
	for i := fromX; i < toX; i++ {
		for j := fromY; j < toY; j++ {
			if comp(f.matrix[i][j]) {
				o := NewObjLocator(FieldPoint{i, j}, f.matrix[i][j])
				result = append(result, o)
			}
		}
	}
	return result
}

func (f *Field) MoveFromTo(from FieldPoint, to FieldPoint) bool {
	if !f.checkCorectPoint(from) || !f.checkCorectPoint(to) {
		return false
	}

	if f.matrix[to.x][to.y].GetType() == Empty {
		f.matrix[to.x][to.y] = f.matrix[from.x][from.y]
		f.matrix[from.x][from.y] = NewEmptyPlace()
		return true
	}

	return false
}

func (f *Field) RemoveFrom(from FieldPoint) {
	if !f.checkCorectPoint(from) {
		return
	}
	f.matrix[from.x][from.y] = NewEmptyPlace()
}

func (f *Field) ClearField() {
	for i := range f.matrix {
		for j := range f.matrix[i] {
			f.matrix[i][j] = NewEmptyPlace()
		}
	}
}

func (f *Field) OnTick() {
	for i := range f.matrix {
		for j := range f.matrix[i] {
			if f.matrix[i][j].GetType() > Doable {
				f.matrix[i][j].(DoableObject).Do(f, FieldPoint{i, j})
			}
		}
	}
}

func (f *Field) Print() {

	maping := func(ot ObjType) string {
		switch ot {
		case Empty:
			return "E"
		case LightSpace:
			return "L"

		case ZooPlanktonT:
			return "Z"
		case PhytoPlanktonT:
			return "P"
		default:
			return "Nan"
		}
	}
	fmt.Println("-------------")
	for i := range f.matrix {
		for j := range f.matrix[i] {
			fmt.Print(maping(f.matrix[i][j].GetType()))
		}
		fmt.Println("")
	}
	fmt.Println("-------------")
}
