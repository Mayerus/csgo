package collections

import (
	"cmp"
	"math/big"
	"testing"

	//"golang.org/x/exp/rand"
	"crypto/rand"
)

type deletion struct {
	value int
	tree  string
}

func RandInt(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}

func CheckAVL[T cmp.Ordered](avl *AvlTree[T], list *LinkedList[T], printList LinkedList[*deletion], t *testing.T) {

	if !avl.containsAllElements(list) {
		for i := 0; i < printList.Count(); i++ {
			record, _ := printList.Get(i)
			t.Logf("Try deleting: %v\n\n%v\n\n\n", record.value, record.tree)
		}
		t.Fatalf("Some values are missing from the AVL tree: \n\n%v\n\nInsert order:%v\n\n\n", avl.String(), list.String())
	}

	if avl.hasDuplicateValues() {
		t.Fatalf("AVL tree has duplicate values: \n\n%v\n\nInsert order:\n\t%v\n\n\n", avl.String(), list.String())
	}

	if !avl.root.isOrdered() {
		t.Fatalf("AVL tree is unordered: \n\n%v\n\nInsert order:\n\t%v\n\n\n", avl.String(), list.String())
	}

	if proper, node := avl.root.propperDynasty(); !proper {
		t.Fatalf("AVL tree has unpropper dynasty: \n\n%v\n\nInsert order:%v\n\nImpaired node:\n\n%v\n\n\n", avl.String(), list.String(), node.String())
	}

	if balanced, discrepancies := avl.root.isBalanced(); !balanced {
		for i := 0; i < printList.Count(); i++ {
			record, _ := printList.Get(i)
			t.Logf("Try deleting: %v\n\n%v\n\n\n", record.value, record.tree)
		}
		t.Fatalf("AVL tree is unbalanced: \n\n%v\n\nInsert order:\n\t%v\n\nDiscrepancies:\n\t%v\n\n\n", avl.String(), list.String(), discrepancies.String())
	}
}

func TestAvlTree(t *testing.T) {
	avl := &AvlTree[int]{}
	list := &LinkedList[int]{}

	k := 1
	l := 20000
	m := 250000
	for i := 0; i < k; i++ {
		for j := 0; j < l; j++ {
			value, err := RandInt(m)
			if err != nil {
				t.Fatalf("%v", err.Error())
				return
			}
			avl.Insert(value)
			if list.Contains(value) == -1 {
				list.Add(value)
			}
		}
	}

	//t.Logf("\n %v\n%v\n\n", avl.String(), list.String())

	// print maximum of last 3 iterations
	printList := &LinkedList[*deletion]{}

	for list.Count() > 0 {
		i, err := RandInt(list.Count())
		if err != nil {
			t.Fatal(err)
		}
		value, _ := list.Get(i)
		list.DeleteAt(i)

		//t.Logf("Try deleting: %v\n\n", value)
		if !avl.Delete(value) {
			for i := 0; i < printList.Count(); i++ {
				record, _ := printList.Get(i)
				t.Logf("Try deleting: %v\n\n%v\n\n\n", record.value, record.tree)
			}
			t.Logf("Try deleting: %v\n\n", value)
			t.Fatalf("Delete: AVL has missing values: \n\n%v\nFailed delete: %v\nList: %v\n\n\n", avl.String(), value, list)
		}
		//t.Logf("\n" + avl.String())

		//printList.Add(&deletion{value, avl.String()})
		//if printList.Count() > 3 {
		//	printList.DeleteAt(0)
		//}

		CheckAVL(avl, list, *printList, t)
	}

	//CheckAVL(avl, list, t)
}

func (t *AvlTree[T]) hasDuplicateValues() bool {
	values := BSTree[T]{}

	return t.root.hasDuplicateValues(&values)
}

func (n *AvlNode[T]) hasDuplicateValues(values *BSTree[T]) bool {
	if n == nil {
		return false
	}
	if values.Search(n.Value) != nil {
		return true
	}
	values.Insert(n.Value)
	return n.Left.hasDuplicateValues(values) || n.Right.hasDuplicateValues(values)
}

func (n *AvlNode[T]) isBalanced() (bool, *LinkedList[T]) {
	list := &LinkedList[T]{}
	return n.isBalancedAuxilary(list), list
}

func (n *AvlNode[T]) isBalancedAuxilary(list *LinkedList[T]) bool {
	if n == nil {
		return true
	}
	if n.Left == nil && n.Right == nil {
		return true
	}
	lHeight := n.Left.getTreeHeight(0)
	rHeight := n.Right.getTreeHeight(0)
	balanceFactor := lHeight - rHeight
	if balanceFactor < -2 || balanceFactor > 2 || n.balanceFactor != int8(balanceFactor) {
		list.Insert(n.Value)
		return false
	}

	return n.Left.isBalancedAuxilary(list) && n.Right.isBalancedAuxilary(list)
}

func (n *AvlNode[T]) getTreeHeight(height int) int {
	if n == nil {
		return height
	}
	height++
	if n.Left == nil && n.Right == nil {
		return height
	}
	lHeight := n.Left.getTreeHeight(height)
	rHeight := n.Right.getTreeHeight(height)
	if lHeight > rHeight {
		return lHeight
	}
	return rHeight
}

func (n *AvlNode[T]) isOrdered() bool {
	if n == nil {
		return true
	}
	if n.Left != nil {
		if n.Value < n.Left.Value {
			return false
		}

		if !n.Left.isOrdered() {
			return false
		}
	}

	if n.Right != nil {
		if n.Value > n.Right.Value {
			return false
		}
		if !n.Right.isOrdered() {
			return false
		}
	}

	return true
}

func (n *AvlNode[T]) propperDynasty() (bool, *AvlNode[T]) {
	if n == nil {
		return true, nil
	}
	if n.Left != nil {
		if n.Left.Parent != n {
			return false, n
		}
		if proper, _ := n.Left.propperDynasty(); !proper {
			return false, n
		}
	}

	if n.Right != nil {
		if n.Right.Parent != n {
			return false, n
		}
		if proper, _ := n.Right.propperDynasty(); !proper {
			return false, n
		}
	}

	return true, nil
}

func (t *AvlTree[T]) containsAllElements(list *LinkedList[T]) bool {
	for node := list.head; node != nil; node = node.next {
		if t.Search(node.data) == nil {
			return false
		}
	}
	return true
}
