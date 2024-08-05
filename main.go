package main

import (
	"fmt"

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
	fmt.Println(avl)
	avl.Insert(2)
	fmt.Println(avl)
	avl.Insert(-1)
	fmt.Println(avl)
	avl.Insert(1)
	fmt.Println(avl)
	avl.Insert(-3)
	fmt.Println(avl)
	avl.Insert(5)
	fmt.Println(avl)
	avl.Insert(3)
	fmt.Println("avl")
	avl.Insert(4)
	avl.Insert(6)
	avl.Insert(7)
	avl.Insert(8)
	avl.Insert(9)
	fmt.Println(avl)

	return
	bst.Delete(5)
	fmt.Println(bst)
	bst.Delete(-3)
	fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
}
