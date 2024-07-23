package main

import (
	"fmt"

	"mayerus/csgo/collections"
)

func main() {
	bst := &collections.BSTree[int]{Root: nil}
	bst.Insert(2)
	bst.Insert(-1)
	bst.Insert(1)
	bst.Insert(-3)
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(4)
	bst.Insert(6)
	fmt.Println(bst)

	bst.Delete(5)
	fmt.Println(bst)
	bst.Delete(-3)
	fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
}
