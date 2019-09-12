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

var _swaggerJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x57\x4b\x6f\xe3\x36\x10\xbe\xe7\x57\x0c\xd8\x1e\x37\x96\x93\x6d\x0f\xcd\xad\xdd\x57\x0d\x74\x01\xa3\x48\x4f\xed\x1e\x18\x69\x2c\x71\x23\x3e\x3a\x43\x5a\x6b\x04\xfe\xef\x05\x29\xd9\x96\x64\x39\x4d\xbc\xc1\xb6\x37\x8b\xf3\xd0\x37\x1f\x67\xbe\x91\x1f\x2e\x00\x44\x6e\x0d\x07\x8d\x2c\x6e\xe0\xcf\x0b\x00\x00\x21\x9d\xab\x55\x2e\xbd\xb2\x26\xfb\xcc\xd6\x88\x0b\x80\x4f\xaf\xa2\xaf\x23\x5b\x84\xfc\x69\xbe\x9c\x57\x38\x48\x5b\x79\xef\xb8\x67\x6f\x64\x59\x22\x89\x1b\x10\xd7\xb3\xb9\x48\x67\xca\xac\xac\xb8\x81\x87\x36\xa0\x40\xce\x49\xb9\x98\x3b\x7a\xdd\x56\x08\x2e\x90\xb3\x8c\x60\x57\xe0\x2b\xc5\x50\xd8\x3c\x68\x34\x3e\x01\x00\xc5\xe0\x2d\x38\xb2\x6b\x55\x20\x14\xb8\xc6\xda\x3a\x24\x06\x69\x40\x19\x56\x65\xe5\x63\x64\x65\x1b\xf0\xf6\x2f\xa3\x8c\x47\x92\xb9\x87\x46\xf9\x0a\x3e\x6e\xd8\x23\xa9\xa0\xe1\x3d\x62\x71\x27\xf3\x7b\xf8\x79\xb9\x48\xb8\x00\x84\x57\xbe\xc6\x88\x62\xca\x58\xab\x1c\x0d\xe3\x1e\x39\x80\x30\x52\x27\xf7\x0f\xcb\xdf\xd6\xaf\x3b\x3f\x00\x11\xa8\x8e\xa7\x89\x8a\x9b\x2c\x6b\x9a\x66\x56\x9a\x30\xb3\x54\x66\x5d\x12\xce\x4a\x57\x5f\xbe\x9e\xcd\x67\x95\xd7\xb5\x48\x81\xdb\xee\x3d\x6b\x24\xee\xc8\x98\xcf\xe6\xb3\xab\x68\x4d\x36\x51\x59\xf6\xf1\xb8\xb6\xb9\xac\xd3\x43\x3a\xbe\x93\x8c\x4b\xe9\xab\x68\xca\xa4\x53\xd9\xfa\xaa\x35\x38\xe9\x2b\x3e\x30\x9d\x95\xca\x57\xe1\xae\x5f\x80\x6b\x33\xee\x9e\x8f\xaf\xe3\x0a\x08\xff\x0e\xc8\x1e\x1c\x12\x68\x65\x82\xc7\x78\x03\xb2\xae\x6d\x83\xc5\xbe\xe6\xa9\x16\xeb\xce\x75\xa8\xbd\x72\x92\x7c\xb6\xb2\xa4\x2f\x0b\xe9\xa5\xd8\xdb\x3f\xf5\x12\x1c\xf5\x5d\x77\x3e\xd5\x7d\xc7\xd1\x1c\xb4\x96\xb4\x89\xa0\xdf\x10\x4a\x8f\x0c\x12\x0c\x36\xf0\x21\x55\x0d\x8a\x39\x60\xdb\x04\x81\x91\x80\xd0\x59\xf2\xfd\x02\x62\x17\xa5\xb7\x2c\x8a\x98\x25\x4f\x59\xda\xe8\x45\x0c\xee\xfb\x3a\x49\x52\xa3\x47\x1a\xa3\x7d\xe8\xfd\x8e\x1d\xb5\x71\xa9\x43\xd8\x93\x32\x65\x2f\x43\xb2\x7e\xb9\x2c\xed\xe5\xae\x89\xfe\x60\xa4\x45\x31\x76\xd9\x59\xc3\xa4\x55\xa5\x4b\x8a\xbc\xbe\x8d\xb4\x8e\xac\xf1\xea\x14\x61\xac\xc6\x53\xc0\x9e\x71\xfb\xea\x45\x20\xbf\xed\x35\xcb\x09\xdc\xc5\x69\x97\xff\x18\xfc\x3b\x2d\x55\x7d\x0a\x36\x4e\x19\x47\x80\x9f\x09\x69\xa5\x6a\x7c\x14\xd0\xfb\x09\x87\x9d\x6d\x2a\xf8\x6c\xfe\x26\xe7\x87\x90\x9d\x8d\xd2\x34\xd0\x03\x00\x71\x3d\x9f\x8f\x8e\x8e\x65\x22\xcd\x07\xb4\x13\x53\x80\x32\xdd\xd0\x8d\x31\xa5\x6d\x21\x8f\xb2\x01\x88\xef\x09\x57\x31\xd1\x77\x59\x81\x2b\x65\x54\x4c\xcc\xd9\x9b\xf1\x08\xfe\xde\x81\x14\x83\xf8\xed\xa9\x9b\x10\x3f\x3c\x01\xfb\x2f\xb2\xd8\x89\xdc\xd7\xe2\x7d\x47\x64\xe9\xd9\x18\xaf\x7f\xfa\x57\x8c\xb7\xd6\x82\x96\x66\xb3\x03\xca\x27\x9b\x4f\xfc\xf8\x94\xeb\x8a\x3b\xd1\xc8\x1a\x18\x69\x8d\x04\x18\x71\x7f\xab\xe2\x2f\xc6\xbf\xb6\x83\xf5\x97\xb9\x38\xb8\xbd\x1d\x55\xe2\x68\x45\xf5\x94\xbe\x7d\x5d\xd1\x7e\x11\x0c\x07\x7e\x2c\xe7\x63\xf3\x57\xf5\xfb\x52\x99\x12\x38\xe4\x39\x32\xaf\xc2\x91\x50\x3c\x97\xba\x98\xee\x0c\xc6\xf6\x1f\x06\xbd\x5c\x87\x75\x7f\x7a\x78\x7a\xe4\x8e\xea\x3a\x19\x13\xd7\x25\x21\xa3\xf1\x71\xad\x1e\x2a\x87\xb2\xbf\x5d\x93\x00\xf4\x95\x7e\x2f\x7e\xf6\xee\x33\xe6\x87\xf9\x8a\xbb\xde\x21\x79\x35\xa2\x5f\xa4\x3c\xe9\xc2\x06\x77\xf2\x88\xac\x0f\x35\x74\xd1\xc5\x1f\x53\xb6\x7f\x75\xf2\x77\x32\xbf\x97\x65\x0a\x69\x2b\x98\xe5\x56\x67\x7a\xf7\x75\x68\xd0\x37\x96\xee\xb3\x55\xf7\x1d\xb8\xff\x31\xfc\x54\x4b\x2d\xff\x08\x9d\xc9\x3e\xa2\x4e\x99\xb2\xc6\x76\xe2\xa2\x50\x4a\x03\xc3\xc1\x39\x9b\x3a\x8d\xcc\x6d\x49\xe7\x51\xf7\xb1\x8b\x7f\x69\xea\xa4\x53\x48\x34\x41\xdc\x13\xfa\x71\xe0\x77\x20\x12\x7e\xbd\xbd\x5d\x76\x1c\x3a\xb9\xa9\xad\x2c\xce\xe7\x2d\xa5\x39\x52\x81\x5d\x16\x49\x24\x37\x43\xd6\x94\x47\x3d\xf6\x7f\x5c\x15\x4f\xab\xf5\xf0\x8b\xa4\x45\xf2\x2d\x2e\x60\x39\x12\xd9\x29\x7d\xeb\xf5\xad\x6b\x9f\xff\x2f\x2d\x3a\xf0\xc0\x2f\x52\xbb\xf6\x6f\x9b\xb3\x7d\x1d\x7d\x21\xfa\xda\x1d\x29\x0e\x92\x7b\xb1\xfd\x27\x00\x00\xff\xff\xaf\x8c\x1b\x72\x53\x0f\x00\x00")

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

	info := bindataFileInfo{name: "swagger.json", size: 3923, mode: os.FileMode(420), modTime: time.Unix(1568287643, 0)}
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
