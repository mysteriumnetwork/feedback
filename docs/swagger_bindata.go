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

var _swaggerJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\x4d\x4f\x1b\x31\x10\xbd\xe7\x57\x8c\xdc\x1e\x81\x45\xa8\xbd\x70\x2c\xd0\x36\x87\x4a\xa8\xa2\xa7\xaa\x07\xb3\x3b\xd9\x18\xd6\x1f\x9d\xf1\x36\x44\x88\xff\x5e\xd9\x1b\xb2\x5e\x67\x13\x41\x40\xb4\x27\x16\xcf\xcc\xf3\x9b\xe7\xf9\xc8\xfd\x04\x40\xf0\x42\xd6\x35\x92\x38\x05\x71\x72\x74\x2c\x0e\xc2\x99\x93\x7e\xce\xe2\x14\x82\x03\x80\x28\x6a\xe5\xe7\xed\xf5\xfa\x20\x78\x58\xf6\xc9\xff\x00\xa2\xb4\x86\x5b\x8d\x21\xec\xe7\xfa\x14\x40\xe8\xb6\xf1\xca\x49\xf2\xc5\xcc\x92\x3e\xac\xa4\x97\x62\x6d\xff\x75\xd0\x03\x38\xb2\x55\x5b\x6e\x02\x48\xe7\x1a\x55\x4a\xaf\xac\x29\x6e\xd8\x9a\xf1\x68\x6e\xb5\x96\xb4\x0c\x69\x9c\x11\x4a\x8f\x0c\x12\x0c\x2e\xe0\x4b\xe4\x0e\x8a\xb9\x45\x58\x28\x3f\x87\x96\x91\x80\xd0\x59\xf2\x22\x41\xb0\x0e\x29\xde\x32\xad\x02\x4a\x19\x51\xba\xe8\x69\x08\x4e\x7d\x9d\x24\xa9\xd1\x23\xe5\x6c\xef\x93\x6f\x00\xe1\x97\x0e\x03\x18\x7b\x52\xa6\x4e\x10\xa2\xf5\xee\xb0\xb6\x87\x46\xea\xe8\xf2\x83\x91\xa6\x55\xee\xf2\x68\x6d\x47\xad\xca\x04\x5b\xd0\xf5\x3c\xc8\x9a\x59\x09\x7f\xb7\x8a\x30\x64\xe3\xa9\xc5\xc4\xf8\x70\xf0\x2a\x94\xcf\x91\x4b\x52\x2e\x68\xb6\x8d\x77\xb5\xdd\xe5\x1f\x93\xbf\xd0\x52\x35\xdb\x68\xe3\x98\x31\x23\xfc\x4c\x4a\x33\xd5\xe0\x4e\x42\x9f\x47\x1c\x1e\x6d\x63\xc1\x7b\xeb\x37\xda\x3f\x84\xec\xac\x61\xe4\x41\x57\x03\x88\x93\xe3\xe3\xec\x08\x86\xcf\x7a\x0a\x22\xf6\x07\x74\x1d\x53\x81\x32\xab\xa6\xcb\x39\x71\x39\x47\x2d\x37\xd0\x00\xc4\x7b\xc2\x59\x00\x7a\x57\x54\x38\x53\x46\x05\x60\x2e\xce\xf2\x16\xfc\xbe\x22\x29\x06\xf1\x0f\xdb\x5e\x42\x7c\x78\x02\xf7\x4f\xb2\x82\x20\x16\xb2\x7f\x29\xdf\x0b\x22\x4b\xcf\xe5\xf8\xf1\x29\xfa\x1a\x8f\x64\x64\x03\x8c\xf4\x07\x09\x30\x5c\xf4\x56\x6c\x27\xf9\x57\xf7\x77\x95\x85\x28\x5c\xe8\xb4\x64\x35\xd4\x98\x6d\x86\x64\x34\x5f\x2a\x53\x43\x57\x6a\x15\x83\xb7\xe0\x86\x6d\x9a\x0f\xe1\xdc\xfc\xa2\x2a\x8d\x97\x73\x5b\x96\xc8\x3c\x6b\x37\xda\xfb\xb9\xfa\x05\xb8\x3d\x64\x9b\xac\xa4\x13\x09\x56\xbf\x6a\xb7\x97\x7c\xa2\x70\x96\xd7\xd6\x98\xb0\xe4\x08\x19\x8d\x0f\xcb\xb0\xcf\x1c\xea\x74\x27\xc6\xb6\x4d\xe7\xf3\x7a\x64\xd9\xeb\x1b\x2c\xfb\xae\x08\x1b\xda\x21\x79\x95\xc9\x2f\x22\x4e\x7c\xb0\xc1\x9b\xec\x18\xc6\xc3\xc9\x37\x5d\xc5\x6f\x4a\xb6\xbe\x3a\xfa\x3b\x59\xde\xca\x3a\x86\x74\x19\x1c\x95\x56\x17\x7a\xc9\x1e\x49\xb5\xda\xa0\x5f\x58\xba\x2d\x66\x88\xd5\xb5\x2c\xfb\x0f\x31\x28\xd8\x58\xf7\x3b\xe4\x8c\xf6\x4c\x3a\x65\xea\x06\xbb\xb6\x0b\xe3\x4d\x1a\x18\x76\xcf\xde\xd2\x69\x64\xee\x52\xda\x4f\xba\x6f\xab\xf8\xd7\x96\x4e\x3a\x85\x44\x23\xc2\x3d\xa1\x1e\x07\x7e\xbd\x90\xf0\xf5\xea\xea\x72\xa5\xa1\x93\xcb\xc6\xca\x6a\x7f\xdd\x22\xcc\xc6\x14\x78\x44\x91\x44\x72\x39\x54\x4d\x79\xd4\xb9\xff\xee\xd1\xb8\x75\xc1\x67\xbf\x23\x3a\x26\x6f\xf1\x00\x97\xd9\xa4\x1d\x9b\x6f\x49\xdd\xba\x7e\xd8\xfe\x0f\x25\x3a\xf0\xc0\x3b\xa9\x5d\x13\xed\xce\xa6\x73\xf4\x95\xe4\xeb\x16\xa5\xe8\x47\xee\xe4\xe1\x6f\x00\x00\x00\xff\xff\xd2\x01\x50\x82\xee\x0c\x00\x00")

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

	info := bindataFileInfo{name: "swagger.json", size: 3310, mode: os.FileMode(420), modTime: time.Unix(1568201685, 0)}
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
