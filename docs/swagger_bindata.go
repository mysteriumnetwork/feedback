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

var _swaggerJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\x4d\x4f\x1b\x31\x10\xbd\xe7\x57\x8c\xdc\x1e\x81\x45\xa8\xbd\x70\x2c\xd0\x36\x87\x4a\xa8\xa2\xa7\xaa\x07\xb3\x3b\xd9\x18\xd6\x1f\x9d\xf1\x36\x44\x88\xff\x5e\xd9\x1b\x12\xaf\xb3\x1b\x41\x40\xb4\x27\x36\x9e\x99\xe7\x37\xcf\xf3\xc1\xfd\x04\x40\xf0\x42\xd6\x35\x92\x38\x05\x71\x72\x74\x2c\x0e\xc2\x99\x93\x7e\xce\xe2\x14\x82\x03\x80\x28\x6a\xe5\xe7\xed\xf5\xfa\x20\x78\x58\xf6\xc9\x6f\x00\x51\x5a\xc3\xad\xc6\x10\xf6\x73\x7d\x0a\x20\x74\xdb\x78\xe5\x24\xf9\x62\x66\x49\x1f\x56\xd2\x4b\xb1\xb6\xff\x3a\xd8\x00\x38\xb2\x55\x5b\x6e\x03\x48\xe7\x1a\x55\x4a\xaf\xac\x29\x6e\xd8\x9a\xe1\x68\x6e\xb5\x96\xb4\x0c\x69\x9c\x11\x4a\x8f\x0c\x12\x0c\x2e\xe0\x4b\xe4\x0e\x8a\xb9\x45\x58\x28\x3f\x87\x96\x91\x80\xd0\x59\xf2\x22\x41\xb0\x0e\x29\xde\x32\xad\x02\x4a\x19\x51\xba\xe8\x69\x08\x4e\x7d\x9d\x24\xa9\xd1\x23\xe5\x6c\xef\x93\x6f\x00\xe1\x97\x0e\x03\x18\x7b\x52\xa6\x4e\x10\xa2\xf5\xee\xb0\xb6\x87\x46\xea\xe8\xf2\x83\x91\xa6\x55\xee\xf2\x68\x6d\x07\xad\xca\x04\x5b\xd0\xf5\x3c\xc8\x9a\x59\x09\x7f\xb7\x8a\x30\x64\xe3\xa9\xc5\xc4\xf8\x70\xf0\x2a\x94\xcf\x91\x4b\x52\x2e\x68\x36\xc6\xbb\x1a\x77\xf9\xc7\xe4\x2f\xb4\x54\xcd\x18\x6d\x1c\x32\x66\x84\x9f\x49\x69\xa6\x1a\xdc\x49\xe8\xf3\x80\xc3\xa3\x6d\x28\x78\x6f\xfd\x06\xfb\x87\x90\x9d\x35\x8c\xdc\xeb\x6a\x00\x71\x72\x7c\x9c\x1d\x41\xff\x59\x4f\x41\xc4\xfe\x80\xae\x63\x2a\x50\x66\xd5\x74\x39\x27\x2e\xe7\xa8\xe5\x16\x1a\x80\x78\x4f\x38\x0b\x40\xef\x8a\x0a\x67\xca\xa8\x00\xcc\xc5\x59\xde\x82\xdf\x57\x24\x45\x2f\xfe\x61\xec\x25\xc4\x87\x27\x70\xff\x24\x2b\x08\x62\x21\xfb\x97\xf2\xbd\x20\xb2\xf4\x5c\x8e\x1f\x9f\xa2\xaf\xf1\x48\x46\x36\xc0\x48\x7f\x90\x00\xc3\x45\x6f\xc5\x76\x92\x7f\x75\x7f\x57\x59\x88\xc2\x85\x4e\x4b\x56\x43\x8d\xd9\x66\x48\x46\x73\x77\x5d\xc5\xe0\x2d\xb8\x7e\x87\xe6\xf3\x37\x37\xbf\xa8\x40\x2f\x95\xa9\x81\xdb\xb2\x44\xe6\x59\xbb\xd5\xd9\xcf\x95\x2e\xc0\xed\xa1\xd8\x64\xa5\x9a\x48\xb0\x36\x5b\x76\xbc\xda\x13\x71\xb3\xbc\x46\x63\xc2\x7e\x23\x64\x34\x3e\xec\xc1\x4d\xe6\x50\xa7\xeb\x30\x76\x6c\x3a\x9a\xd7\xd3\xca\x5e\xdf\x60\xb9\x69\x88\xb0\x9c\x1d\x92\x57\x99\xfc\x22\xe2\xc4\x07\xeb\xbd\xc9\x8e\x39\xdc\x1f\x7a\xd3\x55\xfc\xb6\x64\xeb\xab\xa3\xbf\x93\xe5\xad\xac\x63\x48\x97\xc1\x51\x69\x75\xa1\x97\xec\x91\x54\xab\x0d\xfa\x85\xa5\xdb\x62\x86\x58\x5d\xcb\x72\xf3\x21\x7a\xb5\x1a\x4b\x7e\x87\x9c\xd1\x9e\x49\xa7\x4c\xdd\x60\xd7\x71\x61\xb2\x49\x03\xfd\xc6\xd9\x5b\x3a\x8d\xcc\x5d\x4a\xfb\x49\xf7\x6d\x15\xff\xda\xd2\x49\xa7\x90\x68\x40\xb8\x27\xd4\x63\xcf\x6f\x23\x24\x7c\xbd\xba\xba\x5c\x69\xe8\xe4\xb2\xb1\xb2\xda\x5f\xb7\x08\xb3\x35\x05\x1e\x51\x24\x91\x5c\xf6\x55\x53\x1e\x75\xee\xbf\x7b\x2a\x8e\xee\xf6\xec\x5f\x88\x8e\xc9\x5b\x3c\xc0\x65\x36\x64\x87\xe6\x5b\x52\xb7\xae\xfb\xfd\xbf\x94\x68\xcf\x03\xef\xa4\x76\x4d\xb4\x3b\x9b\xce\xd1\x57\x92\xaf\xdb\x91\x62\x33\x72\x27\x0f\x7f\x03\x00\x00\xff\xff\x9a\x1b\xbd\xb0\xe9\x0c\x00\x00")

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

	info := bindataFileInfo{name: "swagger.json", size: 3305, mode: os.FileMode(420), modTime: time.Unix(1568209061, 0)}
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
