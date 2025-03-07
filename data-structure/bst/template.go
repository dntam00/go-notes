package main

// basic binary tree node
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// binary tree traversal framework
func traverse(root *TreeNode) {
	if root == nil {
		return
	}
	traverse(root.Left)
	traverse(root.Right)
}
