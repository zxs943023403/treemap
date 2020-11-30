package test

import (
	"testing"
	"treemap"
	"fmt"
)

func TestDemo(t *testing.T){
	tree := treemap.CreateNewTreeMap()
	tree.Add("aa","bb")
	tree.Add("cc","dd")
	tree.Add("ee","ff")
	tree.Add(11,22)
	tree.Add(33,44)
	tree.Add(55,66)
	tree.Add('a','b')
	tree.Add('c','d')
	tree.Add('e','f')
	fmt.Println(tree.Keys())
	fmt.Println(tree.Get("cc"))
	fmt.Println(tree.Get(33))
	fmt.Println(tree.Get('a'))
	fmt.Println(tree.Get("ff"))
	fmt.Println(tree.Get(44))
	fmt.Println(tree.Get('b'))
}
