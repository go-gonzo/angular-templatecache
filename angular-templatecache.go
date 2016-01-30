package ngcache

import (
	"bytes"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/omeid/gonzo"
	"github.com/omeid/gonzo/context"
)

type Config struct {
	Name   string //Name of the generated file.
	Module string //Name of AngularJS module.
}

func readall(r io.Reader) (string, error) {
	var buff bytes.Buffer
	_, err := buff.ReadFrom(r)
	return buff.String(), err
}

var cacheTemplate = template.Must(template.New("").Funcs(template.FuncMap{"readall": readall}).Parse(`angular.module("{{ .Module }}", []).run(["$templateCache", function($templateCache) { {{ range $path, $file := .Files }} 
$templateCache.put("{{ $path }}", "{{ readall $file | js }}"); 
{{ end }}
}])`))

type cache struct {
	Config
	Files map[string]gonzo.File
}

func Compile(conf Config) gonzo.Stage {
	return func(ctx context.Context, in <-chan gonzo.File, out chan<- gonzo.File) error {

		b := cache{conf, make(map[string]gonzo.File)}

		for file := range in {
			path := file.FileInfo().Name();
			ctx.Infof("Adding %s", path)
			b.Files[path] = file
			defer file.Close() //Close files AFTER we have build our package.
		}

		buff := new(bytes.Buffer)
		err := cacheTemplate.Execute(buff, b)
		if err != nil {
			ctx.Error(err)
			return err
		}

		fi := gonzo.NewFileInfo()
		fi.SetName(b.Name)
		fi.SetSize(int64(buff.Len()))
		sf := gonzo.NewFile(ioutil.NopCloser(buff), fi)

		out <- sf
		return nil
	}
}
