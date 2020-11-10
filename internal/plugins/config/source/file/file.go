// Package file is a file source. Expected format is json
package file

import (
	"encoding/json"
	"github.com/pubgo/xerror"
	"io/ioutil"
	"os"

	"github.com/asim/nitro/v3/config/source"
)

type file struct {
	path string
	data []byte
	opts source.Options
}

var (
	DefaultPath = "config.json"
)

func (f *file) Read() (_ *source.ChangeSet, err error) {
	defer xerror.RespErr(&err)

	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	info, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	xerror.Panic(f.opts.Encoder.Decode(b, &data))
	b = xerror.PanicBytes(json.Marshal(data))

	cs := &source.ChangeSet{
		Format:    format(f.path, f.opts.Encoder),
		Source:    f.String(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (f *file) String() string {
	return "file"
}

func (f *file) Watch() (source.Watcher, error) {
	if _, err := os.Stat(f.path); err != nil {
		return nil, err
	}
	return newWatcher(f)
}

func (f *file) Write(cs *source.ChangeSet) error {
	return nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	path := DefaultPath
	f, ok := options.Context.Value(filePathKey{}).(string)
	if ok {
		path = f
	}
	return &file{opts: options, path: path}
}
