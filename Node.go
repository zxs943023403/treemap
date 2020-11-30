package treemap

import (
	"errors"
	"crypto/md5"
	"bytes"
	"encoding/gob"
	"sync"
	"math/big"
	"fmt"
)

const (
	// RED 红树设为true
	RED bool = true
	// BLACK 黑树设为false
	BLACK bool = false
)

var h = md5.New()
var hm sync.Mutex

type RBNode struct{
	key                 uint64
	inkey               interface{}
	value               interface{}
	color               bool
	left, right, parent *RBNode
}

// rotate() true左旋/false右旋
// 若有根节点变动则返回根节点
func (rbnode *RBNode) rotate(isRotateLeft bool) (*RBNode, error) {
	var root *RBNode
	if rbnode == nil {
		return root, nil
	}
	if !isRotateLeft && rbnode.left == nil {
		return root, errors.New("右旋左节点不能为空")
	} else if isRotateLeft && rbnode.right == nil {
		return root, errors.New("左旋右节点不能为空")
	}
	parent := rbnode.parent
	var isleft bool
	if parent != nil {
		isleft = parent.left == rbnode
	}
	if isRotateLeft {
		grandson := rbnode.right.left
		rbnode.right.left = rbnode
		rbnode.parent = rbnode.right
		rbnode.right = grandson
	} else {
		grandson := rbnode.left.right
		rbnode.left.right = rbnode
		rbnode.parent = rbnode.left
		rbnode.left = grandson
	}
	// 判断是否换了根节点
	if parent == nil {
		rbnode.parent.parent = nil
		root = rbnode.parent
	} else {
		if isleft {
			parent.left = rbnode.parent
		} else {
			parent.right = rbnode.parent
		}
		rbnode.parent.parent = parent
	}
	return root, nil
}

// getGrandParent() 获取父级节点的父级节点
func (rbnode *RBNode) getGrandParent() *RBNode {
	return rbnode.parent.parent
}

// getSibling() 获取兄弟节点
func (rbnode *RBNode) getSibling() *RBNode {
	p := rbnode.parent
	if p.left == rbnode {
		return p.right
	}else{
		return p.left
	}
}
// GetUncle() 父节点的兄弟节点
func (rbnode *RBNode) getUncle() *RBNode {
	return rbnode.parent.getSibling()
}

func (rbnode *RBNode)getLeftMostChild()*RBNode{
	if rbnode.left == nil {
		return rbnode
	}else{
		return rbnode.getLeftMostChild()
	}
}


func NewRBNode(k,v interface{})*RBNode{
	n := new(RBNode)
	bs,err := intf2keyint(k)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	n.inkey = k
	n.key = bs
	n.value = v
	n.color = RED
	return n;
}

func intf2keyint(key interface{}) (uint64, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return 0, err
	}
	bs := buf.Bytes()
	hm.Lock()
	h.Reset()
	h.Write(bs)
	bs = h.Sum(nil)
	hm.Unlock()
	bi := new(big.Int)
	bi.SetBytes(bs)
	return bi.Uint64(),nil
}
