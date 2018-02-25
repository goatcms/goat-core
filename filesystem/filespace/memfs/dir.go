package memfs

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// Dir is single directory
type Dir struct {
	name     string
	filemode os.FileMode
	time     time.Time
	nodes    []os.FileInfo
}

// Name is a directory name
func (d *Dir) Name() string {
	return d.name
}

// Mode is a unix file/directory mode
func (d *Dir) Mode() os.FileMode {
	return d.filemode
}

// ModTime is modification time
func (d *Dir) ModTime() time.Time {
	return d.time
}

// Sys return native system object
func (d *Dir) Sys() interface{} {
	return nil
}

// Size is length in bytes for regular files; system-dependent for others
func (d *Dir) Size() int64 {
	return int64(len(d.nodes))
}

// IsDir return true if node is a directory
func (d *Dir) IsDir() bool {
	return true
}

// GetNodes return nodes for directory
func (d *Dir) GetNodes() []os.FileInfo {
	return d.nodes
}

// GetNode return single node by name
func (d *Dir) GetNode(nodeName string) (os.FileInfo, error) {
	for _, node := range d.nodes {
		if nodeName == node.Name() {
			return node, nil
		}
	}
	return nil, fmt.Errorf("No find node with name " + nodeName)
}

// AddNode add new node to directory (name must be unique in directory)
func (d *Dir) AddNode(newNode os.FileInfo) error {
	for _, node := range d.nodes {
		if newNode.Name() == node.Name() {
			return fmt.Errorf("node named  " + newNode.Name() + " exists")
		}
	}
	d.nodes = append(d.nodes, newNode)
	return nil
}

// RemoveNodeByName remove a node by name
func (d *Dir) RemoveNodeByName(name string) error {
	for i := 0; i < len(d.nodes); i++ {
		if d.nodes[i].Name() == name {
			d.nodes = append(d.nodes[:i], d.nodes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Con not find node to remove (by name " + name + ")")
}

// Remove remove a node by path
func (d *Dir) Remove(nodePath string, emptyOnly bool) error {
	var currentNode os.FileInfo = d
	if nodePath == "." || nodePath == "./" || nodePath == "" {
		return fmt.Errorf("memfs.Dir.Remove: It is not possible to delete myself")
	}
	pathNodes := strings.Split(path.Clean(nodePath), "/")
	lastDirNode := len(pathNodes) - 1
	for i := 0; i < lastDirNode; i++ {
		nodeName := pathNodes[i]
		if currentNode.IsDir() != true {
			return fmt.Errorf("memfs.Dir.Remove: Node by name %v must be dir to get sub node (path %v )", currentNode.Name(), nodePath)
		}
		var dir = currentNode.(*Dir)
		newNode, err := dir.GetNode(nodeName)
		if err != nil {
			return err
		}
		currentNode = newNode
	}
	if currentNode.IsDir() != true {
		return fmt.Errorf("memfs.Dir.Remove: Node by name %v must be dir to get sub node (path %v )", currentNode.Name(), nodePath)
	}
	currentDir := currentNode.(*Dir)
	removeNodeName := pathNodes[len(pathNodes)-1]
	for key, removedNode := range currentDir.nodes {
		if removedNode.Name() != removeNodeName {
			continue
		}
		if emptyOnly && removedNode.IsDir() {
			removedDir := removedNode.(*Dir)
			if len(removedDir.nodes) != 0 {
				return fmt.Errorf("memfs.Dir.Remove: Prevent remove no empty directory %v", nodePath)
			}
		}
		currentDir.nodes = append(currentDir.nodes[:key], currentDir.nodes[key+1:]...)
		return nil
	}
	return fmt.Errorf("memfs.Dir.Remove: Con not find node to remove (by name %v)", removeNodeName)
}

// GetByPath return node by path
func (d *Dir) GetByPath(nodePath string) (os.FileInfo, error) {
	var currentNode os.FileInfo = d
	if nodePath == "." || nodePath == "./" {
		return currentNode, nil
	}
	pathNodes := strings.Split(path.Clean(nodePath), "/")
	for _, nodeName := range pathNodes {
		if currentNode.IsDir() != true {
			return nil, fmt.Errorf("Node by name " + currentNode.Name() + " must be dir to get sub node (path " + nodePath + " )")
		}
		var dir = currentNode.(*Dir)
		newNode, err := dir.GetNode(nodeName)
		if err != nil {
			return nil, err
		}
		currentNode = newNode
	}
	return currentNode, nil
}

// Copy copy directory and return new directories and files tree
func (d *Dir) Copy() (*Dir, error) {
	var err error
	var nodescopy = make([]os.FileInfo, len(d.nodes))
	for i := 0; i < len(d.nodes); i++ {
		if d.nodes[i].IsDir() {
			var dir = d.nodes[i].(*Dir)
			nodescopy[i], err = dir.Copy()
			if err != nil {
				return nil, err
			}
		} else {
			var file = d.nodes[i].(*File)
			nodescopy[i], err = file.Copy()
			if err != nil {
				return nil, err
			}
		}
	}
	return &Dir{
		name:     d.name,
		filemode: d.filemode,
		time:     d.time,
		nodes:    nodescopy,
	}, nil
}

// MkdirAll crete directories recursive
func (d *Dir) MkdirAll(subPath string, filemode os.FileMode) error {
	pathNodes := strings.Split(path.Clean(subPath), "/")
	currentNode := d
	for i, nodeName := range pathNodes {
		newCurrentNode, err := currentNode.GetNode(nodeName)
		if err == nil {
			// get exist path node
			if !newCurrentNode.IsDir() {
				return fmt.Errorf(nodeName + " exists and is not dir in path " + subPath)
			}
			currentNode = newCurrentNode.(*Dir)
			continue
		}
		//create directories
		var subDir *Dir
		for i2 := len(pathNodes) - 1; i2 >= i; i2-- {
			newSubDir := &Dir{
				name:     pathNodes[i2],
				filemode: filemode,
				time:     time.Now(),
				nodes:    []os.FileInfo{},
			}
			if subDir != nil {
				newSubDir.AddNode(subDir)
			}
			subDir = newSubDir
		}
		currentNode.AddNode(subDir)
		return nil
	}
	return nil
}

// ReadFile read file by path
func (d *Dir) ReadFile(subPath string) ([]byte, error) {
	node, err := d.GetByPath(subPath)
	if err != nil {
		return nil, err
	}
	if node.IsDir() {
		return nil, fmt.Errorf("Use ReadFile on directory ")
	}
	var fileNode = node.(*File)
	return fileNode.GetData(), nil
}

// WriteFile write file by path
func (d *Dir) WriteFile(subPath string, data []byte, perm os.FileMode) error {
	dirPath := path.Dir(subPath)
	d.MkdirAll(dirPath, perm)
	node, err := d.GetByPath(subPath)
	if err != nil {
		//create new file if not exist
		node, err = d.GetByPath(dirPath)
		if err != nil {
			return err
		}
		if !node.IsDir() {
			return fmt.Errorf("There is a file on path " + dirPath)
		}
		var baseDir = node.(*Dir)
		baseDir.AddNode(&File{
			name:     path.Base(subPath),
			filemode: perm,
			time:     time.Now(),
			data:     data,
		})
		return nil
	}
	//overwrite file
	if node.IsDir() {
		return fmt.Errorf("Use WriteFile on directory")
	}
	var file = node.(*File)
	file.SetData(data)
	return nil
}
