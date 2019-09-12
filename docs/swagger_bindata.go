// Code generated for package docs by go-bindata DO NOT EDIT. (@generated)
// sources:
// swagger.json
package docs

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _swaggerJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\x4d\x53\xf3\x36\x10\xbe\xe7\x57\xec\xa8\x3d\x02\xa6\x4c\x7b\x28\xc7\x02\x6d\x73\xe8\x0c\xd3\xa1\xa7\x4e\x0f\xc2\xde\x38\x02\xeb\xa3\xbb\x72\x43\x86\xe1\xbf\x77\x24\x3b\x89\xac\xd8\x79\x21\x30\xbc\xef\x09\x47\xfb\xa1\x67\x1f\xed\x3e\xcb\xf3\x0c\x40\xf0\x4a\xd6\x35\x92\xb8\x04\x71\x71\x76\x2e\x4e\xc2\x99\x93\x7e\xc9\xe2\x12\x82\x03\x80\x28\x6a\xe5\x97\xed\xfd\xf6\x20\x78\x58\xf6\xc9\x6f\x00\x51\x21\x97\xa4\x9c\x57\xd6\x84\x5c\x3f\x00\xe1\xbf\x2d\xb2\x07\x87\x04\x5a\x99\xd6\x23\x28\x06\xd9\x34\x76\x85\x55\xbc\xa7\x8f\x2c\xad\xe1\x56\x63\xb8\xf0\xef\xed\x29\x80\xd0\x6d\xe3\x95\x93\xe4\x8b\x85\x25\x7d\x5a\x49\x2f\xc5\xd6\xfe\x4f\x92\xc0\x91\xad\xda\x72\x3f\x81\x74\xae\x51\xa5\x0c\x90\x8a\x07\xb6\x66\x3c\x9a\x5b\xad\x25\xad\x03\xe8\x2b\x42\xe9\x91\x41\x82\xc1\x15\xfc\x16\xab\x06\xc5\xdc\x22\xac\x94\x5f\x42\xcb\x48\x40\xe8\x2c\xf9\xb4\x00\xeb\x90\xe2\x2d\xf3\x2a\x64\x29\x63\x96\x2e\x7a\x1e\x82\x53\x5f\x27\x49\x6a\xf4\x48\x39\xda\xe7\xe4\x1b\x40\xf8\xb5\xc3\x90\x8c\x3d\x29\x53\x27\x19\xa2\xf5\xe9\xb4\xb6\xa7\x46\xea\xe8\xf2\x17\x23\xcd\xab\xdc\x65\x63\x6d\x47\xad\x2a\x3e\x52\xe0\xf5\x3a\xd0\x9a\x59\xc3\xd3\x29\xc2\x50\x8d\xa7\x16\x13\xe3\xcb\xc9\x87\x40\xbe\x4e\x9a\x65\x02\x77\x35\xed\xf2\x95\xc1\xdf\x68\xa9\x9a\x29\xd8\x38\x66\xcc\x00\xbf\x11\xd2\x42\x35\x78\x10\xd0\xaf\x23\x0e\x1b\xdb\x58\xf0\xd1\xfc\x8d\xce\x0f\x21\x3b\x6b\x18\x79\xa0\x07\x00\xe2\xe2\xfc\x3c\x3b\xda\x97\x89\x38\x1f\xd0\x4d\x4c\x05\xca\xf4\x43\x97\x63\xe2\x72\x89\x5a\xee\x65\x03\x10\xdf\x13\x2e\x42\xa2\xef\x8a\x0a\x17\xca\xa8\x90\x98\x8b\xab\x7c\x04\xff\xec\x41\x8a\x41\xfc\xcb\xd4\x4b\x88\x1f\x5f\x81\xfd\x17\x59\x6d\x44\xee\xbd\x78\x6f\x88\x2c\xbd\x19\xe3\xc5\xcf\x5f\xc4\x78\x67\x2d\x68\x69\xd6\x1b\xa0\x3c\xd9\x7c\xe2\xa7\xd7\x3c\x97\xf1\x48\x46\x36\xc0\x48\xff\x21\x01\x06\xdc\x9f\x55\xfc\x2c\xff\xea\xfe\xf6\x55\x88\xc2\x85\xc1\x4d\x76\x54\x8d\xd9\x8a\x4a\x94\xbe\xbb\xae\x62\xf0\x16\xdc\x70\xe0\x73\x39\xcf\xcd\xef\xea\xf7\x5b\x65\x6a\xe0\xb6\x2c\x91\x79\xd1\xee\x09\xc5\x5b\xa9\x0b\xe9\x8e\x60\x6c\xd6\xb3\x26\x92\x5c\xbb\x75\x3f\x3d\x3c\x09\xb9\x59\x5d\x93\x31\x61\x5d\x12\x32\x1a\x1f\xd6\xea\xae\x72\xa8\xd3\xed\x1a\x05\x20\x55\xfa\xad\xf8\xd9\xfb\x07\x2c\x77\xf3\x15\x76\xbd\x43\xf2\x2a\xa3\x5f\xc4\x3c\xf1\xc1\x06\x6f\x72\x40\xd6\x87\x1a\x3a\xef\xe3\xf7\x29\xdb\x5e\x1d\xfd\x9d\x2c\x1f\x65\x1d\x43\xba\x0a\xce\x4a\xab\x0b\xbd\x66\x8f\xa4\x5a\x6d\xd0\xaf\x2c\x3d\x16\x0b\xc4\xea\x5e\x96\xbb\x0f\x31\xe8\xd5\xd8\xf2\x07\xe8\x8c\xf6\x8c\x3a\x65\xea\x06\xbb\x89\x0b\x42\x29\x0d\x0c\x07\xe7\x68\xea\x34\x32\x77\x25\x1d\x47\xdd\x1f\x7d\xfc\x47\x53\x27\x9d\x42\xa2\x11\xe2\x5e\xd1\x8f\x03\xbf\x1d\x91\xf0\xfb\xdd\xdd\x6d\xcf\xa1\x93\xeb\xc6\xca\xea\x78\xde\x62\x9a\x3d\x15\xd8\x64\x91\x44\x72\x3d\x64\x4d\x79\xd4\xb9\xff\x61\x55\x9c\x56\xeb\xe1\x7f\x24\x1d\x92\xcf\x78\x80\xdb\x4c\x64\xc7\xf4\x2d\xe9\x5b\xd7\xfd\xfe\x56\x5a\x74\xe0\x81\x4f\x52\xbb\x26\xda\x9d\x4d\x75\xf4\x83\xe8\xeb\x76\xa4\xd8\x49\xee\xec\xe5\xff\x00\x00\x00\xff\xff\x48\xda\x20\x60\x72\x0d\x00\x00")

func swaggerJsonBytes() ([]byte, error) {
	return bindataRead(
		_swaggerJson,
		"swagger.json",
	)
}

func swaggerJson() (*asset, error) {
	bytes, err := swaggerJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "swagger.json", size: 3442, mode: os.FileMode(420), modTime: time.Unix(1568280794, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"swagger.json": swaggerJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"swagger.json": &bintree{swaggerJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
