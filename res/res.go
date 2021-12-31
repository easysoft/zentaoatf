// Code generated for package res by go-bindata DO NOT EDIT. (@generated)
// sources:
// res/en/messages_en.json
// res/en/usage.txt
// res/zh/messages_zh.json
// res/zh/usage.txt
package res

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

var _resEnMessages_enJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x1c\x8c\xc1\x0a\x02\x31\x0c\x44\xef\xf9\x8a\x21\xe0\x4d\xf7\x03\xf6\x37\xc4\x93\x88\x04\x8c\x61\xa1\x9b\x95\xa6\x82\x50\xfa\xef\xd2\x5c\xe6\x30\x6f\xe6\x75\x02\xb8\x88\xdb\x57\x4c\x79\x05\xab\x5f\x6e\x57\x3e\xcf\x7a\xd7\x08\x31\x0d\x5e\x71\x27\x00\xe8\x99\x00\x6f\xaf\x39\x7d\x1f\xf5\xa9\x3f\xd9\x3f\x45\xf3\x90\xa8\x55\xf1\x28\xd2\xb6\xc3\x53\xb7\xd8\x82\x53\x70\xe2\x41\xc0\x83\xc6\x3f\x00\x00\xff\xff\xd5\x4e\x64\x2a\x74\x00\x00\x00")

func resEnMessages_enJsonBytes() ([]byte, error) {
	return bindataRead(
		_resEnMessages_enJson,
		"res/en/messages_en.json",
	)
}

func resEnMessages_enJson() (*asset, error) {
	bytes, err := resEnMessages_enJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/en/messages_en.json", size: 116, mode: os.FileMode(420), modTime: time.Unix(1589595844, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resEnUsageTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func resEnUsageTxtBytes() ([]byte, error) {
	return bindataRead(
		_resEnUsageTxt,
		"res/en/usage.txt",
	)
}

func resEnUsageTxt() (*asset, error) {
	bytes, err := resEnUsageTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/en/usage.txt", size: 0, mode: os.FileMode(420), modTime: time.Unix(1602322158, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resZhMessages_zhJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\xca\x49\xcc\x4b\x2f\x4d\x4c\x4f\x55\xb2\x52\x50\xaa\xca\x50\xd2\x01\x89\xe5\xa6\x16\x17\x27\xa6\xa7\x16\x2b\x59\x29\x44\x73\x29\x28\x28\x28\x54\x83\x49\x05\x05\xa5\xcc\x14\x90\xba\xb4\xfc\xa2\xf8\xd4\x8a\xc4\xdc\x82\x9c\x54\xb0\x06\xb0\x54\x49\x51\x62\x5e\x71\x4e\x62\x49\x66\x7e\x1e\x48\xcd\x93\x7d\xdd\x4f\x97\x35\x29\xa8\x16\x2b\x81\x15\xd4\x72\x29\x28\xc4\x72\xd5\x02\x02\x00\x00\xff\xff\x17\x39\x6b\xeb\x73\x00\x00\x00")

func resZhMessages_zhJsonBytes() ([]byte, error) {
	return bindataRead(
		_resZhMessages_zhJson,
		"res/zh/messages_zh.json",
	)
}

func resZhMessages_zhJson() (*asset, error) {
	bytes, err := resZhMessages_zhJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/zh/messages_zh.json", size: 115, mode: os.FileMode(420), modTime: time.Unix(1589595849, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _resZhUsageTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func resZhUsageTxtBytes() ([]byte, error) {
	return bindataRead(
		_resZhUsageTxt,
		"res/zh/usage.txt",
	)
}

func resZhUsageTxt() (*asset, error) {
	bytes, err := resZhUsageTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "res/zh/usage.txt", size: 0, mode: os.FileMode(420), modTime: time.Unix(1602322158, 0)}
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
	"res/en/messages_en.json": resEnMessages_enJson,
	"res/en/usage.txt":        resEnUsageTxt,
	"res/zh/messages_zh.json": resZhMessages_zhJson,
	"res/zh/usage.txt":        resZhUsageTxt,
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
	"res": &bintree{nil, map[string]*bintree{
		"en": &bintree{nil, map[string]*bintree{
			"messages_en.json": &bintree{resEnMessages_enJson, map[string]*bintree{}},
			"usage.txt":        &bintree{resEnUsageTxt, map[string]*bintree{}},
		}},
		"zh": &bintree{nil, map[string]*bintree{
			"messages_zh.json": &bintree{resZhMessages_zhJson, map[string]*bintree{}},
			"usage.txt":        &bintree{resZhUsageTxt, map[string]*bintree{}},
		}},
	}},
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
