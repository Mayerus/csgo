package collections

import (
	"fmt"
)

type Grapher[T comparable] interface {
	AddVertex() int
	AddEdge()
	//Neighbors()
}

type Graph[T comparable] struct {
	Counter  int
	Vertices map[int]*Vertex[T]
}

// Weighted graph
type WGraph[T comparable] struct {
	Counter  int
	Vertices map[int]*WVertex[T]
}

type Vertex[T comparable] struct {
	Value T
	Edges map[int]*Vertex[T]
}

// Weighted graph vertex
type WVertex[T comparable] struct {
	Value T
	Edges map[int]*Edge[T]
}

// Weighted graph edge
type Edge[T comparable] struct {
	Weight float64
	Vertex *WVertex[T]
}

func (g *Graph[T]) AddVertex(value T) int {
	g.Counter++
	g.Vertices[g.Counter] = &Vertex[T]{value, map[int]*Vertex[T]{}}
	return g.Counter
}

func (g *WGraph[T]) AddVertex(value T) int {
	g.Counter++
	g.Vertices[g.Counter] = &WVertex[T]{value, map[int]*Edge[T]{}}
	return g.Counter
}

func (g *Graph[T]) AddEdge(idA, idB int) error {
	if _, ok := g.Vertices[idA]; !ok {
		return fmt.Errorf("Vertex %v does not exist", idA)
	}
	if _, ok := g.Vertices[idB]; !ok {
		return fmt.Errorf("Vertex %v does not exist", idB)
	}
	g.Vertices[idA].Edges[idB] = g.Vertices[idB]
	g.Vertices[idB].Edges[idA] = g.Vertices[idA]
	return nil
}

func (g *WGraph[T]) AddEdge(idA, idB int, weight float64) error {
	if _, ok := g.Vertices[idA]; !ok {
		return fmt.Errorf("Vertex %v does not exist", idA)
	}
	if _, ok := g.Vertices[idB]; !ok {
		return fmt.Errorf("Vertex %v does not exist", idB)
	}
	g.Vertices[idA].Edges[idB] = &Edge[T]{weight, g.Vertices[idB]}
	g.Vertices[idB].Edges[idA] = &Edge[T]{weight, g.Vertices[idA]}
	return nil
}

//func (g *Graph[T]) Neighbors(id int) map[int]T {
//	result := map[int]T{}
//
//	return g.Vertices[id].
//	for nId, edge := range g.Vertices[id].Edges {
//		result[nId] =
//		result = append(result, )
//	}
//}
//
//func (g *WGraph[T]) Neighbors(id int) {
//	result := []T{}
//
//	for nId, edge := range g.Vertices[id].Edges{
//		result = append(result, )
//	}
//}
