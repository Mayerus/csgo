package main

import (
	"fmt"
	"mayerus/csgo/collections"
)

func test(p1, p2 *collections.AvlTree[int]) {
	p3 := *p1
	*p1 = *p2
	*p2 = p3
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
	//list := &collections.LinkedList[int]{}
	values := []int{28, 32}
	delValues := []int{28}
	//k := 1
	//l := 20
	//m := 20
	//for i := 0; i < k; i++ {
	//	for j := 0; j < l; j++ {
	//		value := rand.Intn(m) - m/2
	//		avl.Insert(value)
	//		list.Add(value)
	//	}
	//}

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
	//for list.Count() > 0 {
	//	i := rand.Intn(list.Count())
	//	value, _ := list.Get(i)
	//	list.DeleteAt(value)
	//	avl.Delete(value)
	//	fmt.Println(avl)
	//}

	return
	//bst.Delete(5)
	//fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
	//bst.Delete(-3)
	//fmt.Println(bst)
}
