package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"mayerus/csgo/collections"
)

func RandInt(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}

func propperDynasty(n *collections.AvlNode[int]) (bool, *collections.AvlNode[int]) {
	if n == nil {
		return true, nil
	}
	if n.Left != nil {
		if n.Left.Parent != n {
			return false, n
		}
		if proper, _ := propperDynasty(n.Left); !proper {
			return false, n
		}
	}

	if n.Right != nil {
		if n.Right.Parent != n {
			return false, n
		}
		if proper, _ := propperDynasty(n.Right); !proper {
			return false, n
		}
	}

	return true, nil
}

func main() {
	//t1, t2 := &collections.AvlTree[int]{}, &collections.AvlTree[int]{}
	//t1.Insert(0)
	//t1.Insert(1)
	//t2.Insert(2)
	//t2.Insert(3)
	//fmt.Println(t1)
	//fmt.Println(t2)
	//test(t1, t2)
	//fmt.Println(t1)
	//fmt.Println(t2)
	//return

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
	//fmt.Println(bst)

	avl := &collections.AvlTree[int]{}
	list := &collections.LinkedList[int]{}

	//avl.Insert(1)
	//avl.Insert(2)
	//fmt.Println(avl.Height())
	//return
	l := 20
	m := 50
	for j := 0; j < l; j++ {
		value, err := RandInt(m)
		if err != nil {
			panic(err.Error())
		}
		avl.Insert(value)
		if list.Contains(value) == -1 {
			list.Add(value)
		}
	}

	index, err := RandInt(list.Count())
	if err != nil {
		panic(err)
	}
	wedge, _ := list.Get(index)
	list.DeleteAt(index)

	splitable, t1, t2 := avl.AvlSplit(wedge)

	leftList, rightList := &collections.LinkedList[int]{}, &collections.LinkedList[int]{}

	for i := 0; i < list.Count(); i++ {
		value, _ := list.Get(i)
		if value < wedge {
			leftList.Add(value)
			continue
		}
		if value > wedge {
			rightList.Add(value)
			continue
		}
	}

	fmt.Printf("\table: %v\n\nt1:\n\n%s\n\nt2:\n\n%s\n\n\n", splitable, t1, t2)

	//leftAvl := &collections.AvlTree[int]{}
	//rightAvl := &collections.AvlTree[int]{}
	//n := 2.0
	//a := 2.0
	//smallTreeSize := int(math.Pow(2, n-1) + 1)
	//bigTreeSize := int(math.Pow(2, n-1+a) + 1)
	return

	values := []int{28, 32}
	delValues := []int{28}

	for _, v := range values {
		avl.Insert(v)
	}

	//fmt.Println(list)
	fmt.Println(avl)

	for _, v := range delValues {
		result := avl.Delete(v)
		fmt.Println(result)
		fmt.Println(avl)
	}

	return
	//bst.Delete(5)
	//fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
}
