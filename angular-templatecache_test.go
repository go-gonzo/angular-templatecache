package ngcache

import "testing"



var tests = []struct{
	config Config
	Files map[string]string
	Expect string
}{
	{
		Config{"templaes.js", "bing", "app"},
		map[string]string{
			"test.html": "<test>testing, testing!</test>",
			"hello.html": "<hello>World!</hello>",
		},
		`app.module("bing")`,
	},

}


func TestCompile(t *testing.T) {
}
