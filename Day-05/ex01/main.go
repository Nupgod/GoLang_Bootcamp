package main

import (
	"fmt"
)

type TreeNode struct {
	Let    rune
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func unrollGarland(root *TreeNode) ([]bool, []rune) {
	if root == nil {
		return nil, nil
	}

	var result []bool
	var lets []rune
	currentLevel := []*TreeNode{root}
	nextLevel := []*TreeNode{}
	level := 1 // Start with level 1

	for len(currentLevel) > 0 {
		node := currentLevel[len(currentLevel)-1]
		currentLevel = currentLevel[:len(currentLevel)-1]

		result = append(result, node.HasToy)
		lets = append(lets, node.Let)

		// Check if the current level is even or odd
		if level%2 == 0 { // Even level, traverse from left to right
			if node.Left != nil {
				nextLevel = append(nextLevel, node.Left)
			}
			if node.Right != nil {
				nextLevel = append(nextLevel, node.Right)
			}
		} else { // Odd level, traverse from right to left
			if node.Right != nil {
				nextLevel = append(nextLevel, node.Right)
			}
			if node.Left != nil {
				nextLevel = append(nextLevel, node.Left)
			}
		}

		if len(currentLevel) == 0 {
			currentLevel, nextLevel = nextLevel, currentLevel
			level++ // Increment the level for the next iteration
		}
	}

	return result, lets
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
		fmt.Printf("1: %c\n", root.Let)
	} else {
		fmt.Printf("0: %c\n", root.Let)
	}

	printTree(root.Left, level+1)
}

func main() {
	root := &TreeNode{
		Let:    'A',
		HasToy: true,
		Left: &TreeNode{
			Let:    'C',
			HasToy: true,
			Left: &TreeNode{
				Let:    'D',
				HasToy: true,
			},
			Right: &TreeNode{
				Let:    'E',
				HasToy: false,
			},
		},
		Right: &TreeNode{
			Let:    'B',
			HasToy: false,
			Left: &TreeNode{
				Let:    'F',
				HasToy: true,
			},
			Right: &TreeNode{
				Let:    'G',
				HasToy: true,
			},
		},
	}
	fmt.Println("Tree:")
	printTree(root, 0)
	ans, lets := unrollGarland(root)
	fmt.Println(ans) // true true false true true false true
	fmt.Printf("%c\n", lets)

	root = &TreeNode{
		Let:    'A',
		HasToy: false,
		Left: &TreeNode{
			Let:    'B',
			HasToy: true,
			Left: &TreeNode{
				Let:    'D',
				HasToy: true,
			},
			Right: &TreeNode{
				Let:    'E',
				HasToy: false,
			},
		},
		Right: &TreeNode{
			Let:    'C',
			HasToy: false,
			Left: &TreeNode{
				Let:    'F',
				HasToy: true,
				Left: &TreeNode{
					Let:    'H',
					HasToy: true,
				},
				Right: &TreeNode{
					Let:    'I',
					HasToy: true,
				},
			},
			Right: &TreeNode{
				Let:    'G',
				HasToy: true,
				Left: &TreeNode{
					Let:    'F',
					HasToy: true,
					},
				Right: &TreeNode{
						Let:    'I',
						HasToy: true,
					},
				
			},
		},
	}
	fmt.Println("Tree:")
	printTree(root, 0)
	ans, lets = unrollGarland(root)
	fmt.Println(ans) // true true false true true false true
	fmt.Printf("%c\n", lets)
}


