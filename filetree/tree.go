package filetree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// supportedExts is the set of file extensions shown in the sidebar.
var supportedExts = map[string]bool{
	".md":   true,
	".txt":  true,
	".text": true,
	".rst":  true,
	".log":  true,
	"":      false, // no extension — handled separately
}

// Node represents a file or directory in the tree.
type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*Node
	Expanded bool
	Depth    int
}

// IsSupportedFile returns true for directories and supported text file extensions.
func IsSupportedFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		return false
	}
	return supportedExts[ext]
}

// ExpandNode loads children for a directory node that hasn't been expanded yet.
func ExpandNode(node *Node) error {
	if !node.IsDir || node.Expanded {
		return nil
	}
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	node.Children = nil
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if !e.IsDir() && !IsSupportedFile(name) {
			continue
		}
		node.Children = append(node.Children, &Node{
			Name:  name,
			Path:  filepath.Join(node.Path, name),
			IsDir: e.IsDir(),
			Depth: node.Depth + 1,
		})
	}
	node.Expanded = true
	return nil
}

// Flatten returns a depth-first ordered slice of visible nodes (respects Expanded).
func Flatten(root *Node) []*Node {
	return flatten(root)
}

func flatten(node *Node) []*Node {
	var result []*Node
	// Don't include the root node itself — show its children at top level
	if node.Depth == 0 {
		if node.Expanded {
			for _, child := range node.Children {
				result = append(result, flattenNode(child)...)
			}
		}
		return result
	}
	return flattenNode(node)
}

func flattenNode(node *Node) []*Node {
	result := []*Node{node}
	if node.IsDir && node.Expanded {
		for _, child := range node.Children {
			result = append(result, flattenNode(child)...)
		}
	}
	return result
}
