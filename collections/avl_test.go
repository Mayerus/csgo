package collections

import (
	"cmp"
	"math"
	"math/big"
	"testing"

	//"golang.org/x/exp/rand"
	"crypto/rand"
)

const (
	MaxElements = 20
	MaxValue    = 250000
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

func CheckAVL[T cmp.Ordered](avl *AvlTree[T], list *LinkedList[T], printList *LinkedList[*deletion], t *testing.T) {

	if !avl.containsAllElements(list) {
		if printList != nil {
			for i := 0; i < printList.Count(); i++ {
				record, _ := printList.Get(i)
				t.Logf("Try deleting: %v\n\n%v\n\n\n", record.value, record.tree)
			}
		}
		t.Fatalf("Some values are missing from the AVL tree: \n\n%v\n\nInsert order:%v\n\n\n", avl.String(), list.String())
	}

	if avl.hasDuplicateValues() {
		t.Fatalf("AVL tree has duplicate values: \n\n%v\n\nInsert order:\n\t%v\n\n\n", avl.String(), list.String())
	}

	if !avl.root.isOrdered() {
		t.Fatalf("AVL tree is unordered: \n\n%v\n\nInsert order:\n\t%v\n\n\n", avl.String(), list.String())
	}

	if !avl.root.propperDynasty(t) {
		t.Fatalf("AVL tree has unpropper dynasty: \n\n%v\n\nInsert order:%v\n\n\n", avl.String(), list.String())
	}

	if balanced, discrepancies := avl.root.isBalanced(); !balanced {
		if printList != nil {
			for i := 0; i < printList.Count(); i++ {
				record, _ := printList.Get(i)
				t.Logf("Try deleting: %v\n\n%v\n\n\n", record.value, record.tree)
			}
		}
		t.Fatalf("AVL tree is unbalanced: \n\n%v\n\nInsert order:\n\t%v\n\nDiscrepancies:\n\t%v\n\n\n", avl.String(), list.String(), discrepancies.String())
	}
}

func TestAvlInsertion(t *testing.T) {
	avl := &AvlTree[int]{}
	list := &LinkedList[int]{}

	k := 1
	for i := 0; i < k; i++ {
		for j := 0; j < MaxElements; j++ {
			value, err := RandInt(MaxValue)
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

	CheckAVL(avl, list, nil, t)
}

func TestAvlDeletion(t *testing.T) {
	avl := &AvlTree[int]{}
	list := &LinkedList[int]{}

	k := 1
	for i := 0; i < k; i++ {
		for j := 0; j < MaxElements; j++ {
			value, err := RandInt(MaxValue)
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

		CheckAVL(avl, list, printList, t)
	}
}

func TestAvlRightJoin(t *testing.T) {
	leftAvl := &AvlTree[int]{}
	rightAvl := &AvlTree[int]{}

	n := 5.0
	a := 7.0

	smallTreeSize := int(math.Pow(2, n-1) + 1)
	bigTreeSize := int(math.Pow(2, n-1+a) + 1)

	for i := 0; i < bigTreeSize; i++ {
		leftAvl.Insert(i)
	}
	for i := bigTreeSize + 1; i < bigTreeSize+smallTreeSize+1; i++ {
		rightAvl.Insert(i)
	}
	//t.Log("\n" + leftAvl.String())
	//t.Log("\n" + rightAvl.String())

	joined, joinedTree := AvlJoin(leftAvl, rightAvl, bigTreeSize)

	if !joined {
		t.Fatalf("AVL right join failed: \n\nLeft tree:\n\n%v\n\nRight tree:\n\n%v\n\n Join value: %v\n\n", leftAvl.String(), rightAvl.String(), bigTreeSize)
	}

	list := &LinkedList[int]{}
	for i := 0; i < smallTreeSize+bigTreeSize+1; i++ {
		list.Add(i)
	}

	emptyList := &LinkedList[*deletion]{}

	CheckAVL(joinedTree, list, emptyList, t)
}

func TestAvlLeftJoin(t *testing.T) {
	leftAvl := &AvlTree[int]{}
	rightAvl := &AvlTree[int]{}

	n := 5.0
	a := 7.0

	smallTreeSize := int(math.Pow(2, n-1) + 1)
	bigTreeSize := int(math.Pow(2, n-1+a) + 1)

	for i := 0; i < smallTreeSize; i++ {
		leftAvl.Insert(i)
	}

	for i := smallTreeSize + 1; i < bigTreeSize+smallTreeSize+1; i++ {
		rightAvl.Insert(i)
	}
	//t.Log("\n" + leftAvl.String())
	//t.Log("\n" + rightAvl.String())

	joined, joinedTree := AvlJoin(leftAvl, rightAvl, smallTreeSize)

	if !joined {
		t.Fatalf("AVL left join failed: \n\nLeft tree:\n\n%v\n\nRight tree:\n\n%v\n\n Join value: %v\n\n", leftAvl.String(), rightAvl.String(), smallTreeSize)
	}

	list := &LinkedList[int]{}
	for i := 0; i < smallTreeSize+bigTreeSize+1; i++ {
		list.Add(i)
	}

	emptyList := &LinkedList[*deletion]{}

	CheckAVL(joinedTree, list, emptyList, t)
}

func TestAvlSplit(t *testing.T) {
	avl := &AvlTree[int]{}
	list := &LinkedList[int]{}

	k := 1
	for i := 0; i < k; i++ {
		for j := 0; j < MaxElements; j++ {
			value, err := RandInt(MaxValue)
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

	index, err := RandInt(list.Count())
	if err != nil {
		t.Fatal(err)
	}
	wedge, _ := list.Get(index)
	list.DeleteAt(index)

	splitable, t1, t2 := avl.AvlSplit(wedge)

	leftList, rightList := &LinkedList[int]{}, &LinkedList[int]{}

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

	t.Log("\nSplitable: ", splitable, "\n\nt1:\n\n", t1.String(), "\n\nt2:\n\n", t2.String(), "\n\n\n")

	CheckAVL(t1, leftList, nil, t)
	CheckAVL(t2, rightList, nil, t)
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

func (n *AvlNode[T]) propperDynasty(t *testing.T) bool {
	if n == nil {
		return true
	}
	if n.Left != nil {
		if n.Left.Parent != n {
			t.Log("impropper: \n", n.Left.String(), "\n n: ", &n, "\n n.Left.parent: ", &n.Left.Parent, "\n\n")
			return false
		}
		if !n.Left.propperDynasty(t) {
			return false
		}
	}

	if n.Right != nil {
		if n.Right.Parent != n {
			t.Log("impropper: \n", n.Left.String(), "\n n: ", &n, "\n n.Right.parent: ", &n.Left.Parent)
			return false
		}
		if !n.Right.propperDynasty(t) {
			return false
		}
	}

	return true
}

func (t *AvlTree[T]) containsAllElements(list *LinkedList[T]) bool {
	for node := list.head; node != nil; node = node.next {
		if t.Search(node.data) == nil {
			return false
		}
	}
	return true
}
