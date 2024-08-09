package main

import (
	"fmt"
	"math/rand"
	"mayerus/csgo/collections"
)

func main() {
	bst := &collections.BSTree[int]{}
	bst.Insert(2)
	bst.Insert(-1)
	bst.Insert(1)
	bst.Insert(-3)
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(4)
	bst.Insert(6)
	bst.Insert(7)
	bst.Insert(8)
	bst.Insert(9)
	fmt.Println(bst)

	avl := &collections.AvlTree[int]{}

	k := 1
	l := 100
	m := 100
	for i := 0; i < k; i++ {
		for j := 0; j < l; j++ {
			value := rand.Intn(m) - m/2
			avl.Insert(value)
		}
	}

	fmt.Println(avl)
	return
	//bst.Delete(5)
	//fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
}
