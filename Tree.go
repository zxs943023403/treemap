package treemap

import (
	"fmt"
	"errors"
)

const (
	// 左旋
	LEFTROTATE bool = true
	// 右旋
	RIGHTROTATE bool = false
)

type TreeMap struct {
	root *RBNode
}

func CreateNewTreeMap()*TreeMap {
	tree := new(TreeMap)
	tree.root = NewRBNode("",nil)
	tree.root.color = BLACK
	return tree
}

func (rbtree *TreeMap)Add(k,v interface{}){
	rbtree.insertNode(rbtree.root,k,v)
}

func (rbtree *TreeMap)Get(k interface{})(interface{},error){
	ik,err := intf2keyint(k)
	if err != nil {
		return nil,err
	}
	sm := get(rbtree.root,ik)
	if sm == nil {
		return nil,errors.New("不存在")
	}
	return sm.value,nil
}

func (rbtree *TreeMap)Keys()[]interface{}{
	ks := make([]interface{},0,0)
	return keys(rbtree.root,ks)
}

func keys(n *RBNode,ks []interface{})[]interface{}{
	if n == nil {
		return ks
	}
	ks = keys(n.left,ks)
	ks = append(ks,n.inkey)
	ks = keys(n.right,ks)
	return ks
}

func get(n *RBNode,k uint64)*RBNode{
	if n == nil {
		return nil
	}
	if n.key > k {
		return get(n.left,k)
	}else if n.key < k {
		return get(n.right,k)
	}else{
		return n
	}
}

func (rbtree *TreeMap) insertNode(pnode *RBNode, k,v interface{}) {
	in,_ := intf2keyint(k)
	if pnode.key == in {
		return
	}
	if pnode.key > in {
		// 插入数据不大于父节点，插入左节点
		if pnode.left != nil {
			rbtree.insertNode(pnode.left, k,v)
		} else {
			tmpnode := NewRBNode(k,v)
			tmpnode.parent = pnode
			pnode.left = tmpnode
			rbtree.insertCheck(tmpnode)
		}
	} else {
		// 插入数据大于父节点，插入右节点
		if pnode.right != nil {
			rbtree.insertNode(pnode.right, k,v)
		} else {
			tmpnode := NewRBNode(k,v)
			tmpnode.parent = pnode
			pnode.right = tmpnode
			rbtree.insertCheck(tmpnode)
		}
	}
}
func (rbtree *TreeMap) insertCheck(node *RBNode) {
	if node.parent == nil {
		// 检查1：若插入节点没有父节点，则该节点为root
		rbtree.root = node
		// 设置根节点为black
		rbtree.root.color = BLACK
		return
	}

	// 父节点是黑色的话直接添加，红色则进行处理
	if node.parent.color == RED {
		if node.getUncle() != nil && node.getUncle().color == RED {
			// 父节点及叔父节点都是红色，则转为黑色
			node.parent.color = BLACK
			node.getUncle().color = BLACK
			// 祖父节点改成红色
			node.getGrandParent().color = RED
			// 递归处理
			rbtree.insertCheck(node.getGrandParent())
		} else {
			// 父节点红色，父节点的兄弟节点不存在或为黑色
			isleft := node == node.parent.left
			isparentleft := node.parent == node.getGrandParent().left
			if !isleft && isparentleft {
				rbtree.rotateLeft(node.parent)
				rbtree.rotateRight(node.parent)

				node.color = BLACK
				node.left.color = RED
				node.right.color = RED
			} else if isleft && !isparentleft {
				rbtree.rotateRight(node.parent)
				rbtree.rotateLeft(node.parent)

				node.color = BLACK
				node.left.color = RED
				node.right.color = RED
			} else if isleft && isparentleft {
				node.parent.color = BLACK
				node.getGrandParent().color = RED
				rbtree.rotateRight(node.getGrandParent())
			} else if !isleft && !isparentleft {
				node.parent.color = BLACK
				node.getGrandParent().color = RED
				rbtree.rotateLeft(node.getGrandParent())
			}
		}
	}
}
func (rbtree *TreeMap) rotateLeft(node *RBNode) {
	if tmproot, err := node.rotate(LEFTROTATE); err == nil {
		if tmproot != nil {
			rbtree.root = tmproot
		}
	} else {
		fmt.Println(err.Error())
	}
}
func (rbtree *TreeMap) rotateRight(node *RBNode) {
	if tmproot, err := node.rotate(RIGHTROTATE); err == nil {
		if tmproot != nil {
			rbtree.root = tmproot
		}
	} else {
		fmt.Println(err.Error())
	}
}

// 删除对外方法
func (rbtree *TreeMap) Delete(data uint64) {
	rbtree.delete_child(rbtree.root, data)
}

// 删除节点
func (rbtree *TreeMap) delete_child(n *RBNode, data uint64) bool {
	if data < n.key {
		if n.left == nil {
			return false
		}
		return rbtree.delete_child(n.left, data)
	}
	if data > n.key {
		if n.right == nil {
			return false
		}
		return rbtree.delete_child(n.right, data)
	}

	if n.right == nil || n.left == nil {
		// 两个都为空或其中一个为空，转为删除一个子树节点的问题
		rbtree.delete_one(n)
		return true
	}

	//两个都不为空，转换成删除只含有一个子节点节点的问题
	mostLeftChild := n.right.getLeftMostChild()
	tmpval := n.key
	n.key = mostLeftChild.key
	mostLeftChild.key = tmpval

	rbtree.delete_one(mostLeftChild)

	return true
}

// 删除只有一个子节点的节点
func (rbtree *TreeMap) delete_one(n *RBNode) {
	var child *RBNode
	isadded := false
	if n.left == nil {
		child = n.right
	} else {
		child = n.left
	}

	if n.parent == nil && n.left == nil && n.right == nil {
		n = nil
		rbtree.root = n
		return
	}
	if n.parent == nil {
		n = nil
		child.parent = nil
		rbtree.root = child
		rbtree.root.color = BLACK
		return
	}

	if n.color == RED {
		if n == n.parent.left {
			n.parent.left = child

		} else {
			n.parent.right = child
		}
		if child != nil {
			child.parent = n.parent
		}
		n = nil
		return
	}

	if child != nil && child.color == RED && n.color == BLACK {
		if n == n.parent.left {
			n.parent.left = child

		} else {
			n.parent.right = child
		}
		child.parent = n.parent
		child.color = BLACK
		n = nil
		return
	}

	// 如果没有孩子节点，则添加一个临时孩子节点
	if child == nil {
		child = NewRBNode(0,nil)
		child.parent = n
		isadded = true
	}

	if n.parent.left == n {
		n.parent.left = child
	} else {
		n.parent.right = child
	}

	child.parent = n.parent

	if n.color == BLACK {
		if !isadded && child.color == RED {
			child.color = BLACK
		} else {
			rbtree.deleteCheck(child)
		}
	}

	// 如果孩子节点是后来加的，需删除
	if isadded {
		if child.parent.left == child {
			child.parent.left = nil
		} else {
			child.parent.right = nil
		}
		child = nil
	}
	n = nil
}

// deleteCheck() 删除验证
func (rbtree *TreeMap) deleteCheck(n *RBNode) {
	if n.parent == nil {
		n.color = BLACK
		return
	}
	if n.getSibling().color == RED {
		n.parent.color = RED
		n.getSibling().color = BLACK
		if n == n.parent.left {
			rbtree.rotateLeft(n.parent)
		} else {
			rbtree.rotateRight(n.parent)
		}
	}
	//注意：这里n的兄弟节点发生了变化，不再是原来的兄弟节点
	is_parent_red := n.parent.color
	is_sib_red := n.getSibling().color
	is_sib_left_red := BLACK
	is_sib_right_red := BLACK
	if n.getSibling().left != nil {
		is_sib_left_red = n.getSibling().left.color
	}
	if n.getSibling().right != nil {
		is_sib_right_red = n.getSibling().right.color
	}
	if !is_parent_red && !is_sib_red && !is_sib_left_red && !is_sib_right_red {
		n.getSibling().color = RED
		rbtree.deleteCheck(n.parent)
		return
	}
	if is_parent_red && !is_sib_red && !is_sib_left_red && !is_sib_right_red {
		n.getSibling().color = RED
		n.parent.color = BLACK
		return
	}
	if n.getSibling().color == BLACK {
		if n.parent.left == n && is_sib_left_red && !is_sib_right_red {
			n.getSibling().color = RED
			n.getSibling().left.color = BLACK
			rbtree.rotateRight(n.getSibling())
		} else if n.parent.right == n && !is_sib_left_red && is_sib_right_red {
			n.getSibling().color = RED
			n.getSibling().right.color = BLACK
			rbtree.rotateLeft(n.getSibling())
		}
	}
	n.getSibling().color = n.parent.color
	n.parent.color = BLACK
	if n.parent.left == n {
		n.getSibling().right.color = BLACK
		rbtree.rotateLeft(n.parent)
	} else {
		n.getSibling().left.color = BLACK
		rbtree.rotateRight(n.parent)
	}
}
