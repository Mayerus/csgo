package main

import (
	"os"

	"mayerus/csgo/collections"
)

func main() {
	bst := &collections.BSTree[int]{Root: nil}
	bst.Insert(1)
	bst.Insert(-1)
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(4)
	bst.Insert(6)
	bst.Print(os.Stdout)

	bst.Delete(4)
	bst.Print(os.Stdout)
	bst.Insert(5)
	bst.Print(os.Stdout)
}
