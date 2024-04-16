package main

import (
	"fmt"
)

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func areToysBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	leftCount := countToys(root.Left)
	rightCount := countToys(root.Right)

	return leftCount == rightCount
}

func countToys(node *TreeNode) int {
	if node == nil {
		return 0
	}

	count := 0
	if node.HasToy {
		count = 1
	}

	return count + countToys(node.Left) + countToys(node.Right)
}


func printTree(root *TreeNode, level int) {
	if root == nil {
		return
	}

	printTree(root.Right, level+1)

	for i := 0; i < level; i++ {
		fmt.Print("    ")
	}

	if root.HasToy {
		fmt.Println("1")
	} else {
		fmt.Println("0")
	}

	printTree(root.Left, level+1)
}


func main() {
	root := &TreeNode{
		HasToy: true,
		Left: &TreeNode{
			HasToy: false,
			Left: &TreeNode{
				HasToy: false,
			},
			Right: &TreeNode{
				HasToy: true,
			},
		},
		Right: &TreeNode{
			HasToy: true,
			Left: &TreeNode{
				HasToy: true,
			},
			Right: &TreeNode{
				HasToy: false,
			},
		},
	}

	fmt.Println("Tree:")
	printTree(root, 0)
	fmt.Println("Are toys balanced?", areToysBalanced(root)) // Output: false

	root = &TreeNode{
		HasToy: false,
		Left: &TreeNode{
			HasToy: false,
			Left: &TreeNode{
				HasToy: false,
			},
			Right: &TreeNode{
				HasToy: true,
			},
		},
		Right: &TreeNode{
			HasToy: true,
			Left: &TreeNode{
				HasToy: true,
			},
			Right: &TreeNode{
				HasToy: false,
			},
		},
	}

	fmt.Println("\nTree:")
	printTree(root, 0)
	fmt.Println("Are toys balanced?", areToysBalanced(root)) // Output: true

	root = &TreeNode{
		HasToy: false,
		Left: &TreeNode{
			HasToy: true,
			Left: &TreeNode{
				HasToy: false,
			},
			Right: &TreeNode{
				HasToy: true,
			},
		},
		Right: &TreeNode{
			HasToy: true,
			Left: &TreeNode{
				HasToy: true,
			},
			Right: &TreeNode{
				HasToy: false,
			},
		},
	}

	fmt.Println("\nTree:")
	printTree(root, 0)
	fmt.Println("Are toys balanced?", areToysBalanced(root)) // Output: true

	root = &TreeNode{
		HasToy: true,
		Left: &TreeNode{
			HasToy: true,
			Left: &TreeNode{
				HasToy: false,
				Left: &TreeNode{
					HasToy: false,
				},
				Right: &TreeNode{
					HasToy: true,
				},
			},
			Right: &TreeNode{
				HasToy: true,
			},
		},
		Right: &TreeNode{
			HasToy: true,
			Left: &TreeNode{
				HasToy: true,
			},
			Right: &TreeNode{
				HasToy: false,
				Right: &TreeNode{
					HasToy: true,
				},
			},
		},
	}

	fmt.Println("\nTree:")
	printTree(root, 0)
	fmt.Println("Are toys balanced?", areToysBalanced(root)) // Output: true
}