package leetcode

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func minDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	q := make([]*TreeNode, 1)
	q[0] = root
	depth := 1
	for len(q) > 0 {
		currentSize := len(q)
		for i := 0; i < currentSize; i++ {
			current := q[i]
			if current.Left == nil && current.Right == nil {
				return depth
			}
			if current.Left != nil {
				q = append(q, current.Left)
			}
			if current.Right != nil {
				q = append(q, current.Right)
			}
		}
		q = q[currentSize:]
		depth++
	}
	return depth
}
