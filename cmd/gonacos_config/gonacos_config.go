package main // import "github.com/reatang/go-simple-nacos/cmd/gonacos_config"

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const TypeStruct = "struct"
const TypeMap = "map"

var (
	configTemplate = `
var __{{ .ConfigName }} = atomic.Value{}

func Register{{ .ConfigName }}(ncc *gonacos.NacosConfigClient) {
	ncc.WatchF("{{ .Group }}", "{{ .DataId }}", func(namespace, group, dataId, data string) {
		var c {{ .ConfigName }}
		{{ .DecodeFuncCode }}
		if err == nil {
			__{{ .ConfigName }}.Store(c)
		}
	})
}

func Get{{ .ConfigName }}() {{ .ConfigName }} {
	return __{{ .ConfigName }}.Load().({{ .ConfigName }})
}
`
)

var (
	configName = flag.String("config", "", "结构体名称 或者 map变量名称")
	dataid     = flag.String("dataid", "", "nacos内的dataid名称")
	group      = flag.String("group", "", "nacos内的group名称，不填则使用NacosConfigClient的defaultGroup")
	codec      = flag.String("codec", "yaml", "编码格式，支持json、yaml")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of gonacos_config:\n")
	fmt.Fprintf(os.Stderr, "\tgonacos_config --config=SomeConfig --group=DEFAULT_GROUP --dataid=config --codec=yaml\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("gonacos_config: ")
	flag.Usage = Usage
	flag.Parse()

	if len(*configName) == 0 || len(*dataid) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	currentPackage := os.Getenv("GOPACKAGE")
	if currentPackage == "" {
		log.Println("请在 go:generate 上下文中使用")
		os.Exit(2)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		os.Exit(2)
	}

	g := Generator{
		configName: *configName,
		dataid:     *dataid,
		group:      *group,
		codec:      *codec,
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"gonacos_config %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("package %s", currentPackage)
	g.Printf("\n")
	g.Printf("import \"sync/atomic\"\n")
	g.Printf("import gonacos \"github.com/reatang/go-simple-nacos\"\n")

	// 生成代码
	err = g.generate()
	if err != nil {
		log.Println("生成失败：" + err.Error())
		return
	}

	// 格式化并生成代码
	src := g.format()

	// Write to file.
	baseName := fmt.Sprintf("nacos_%s.gen.go", strings.ToLower(*dataid))
	outputName := filepath.Join(dir, strings.ToLower(baseName))
	err = os.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("写文件失败: %s", err)
	}
}

type Generator struct {
	buf bytes.Buffer

	configName string
	dataid     string
	group      string
	codec      string
}

func (g *Generator) Printf(format string, args ...interface{}) {
	_, err := fmt.Fprintf(&g.buf, format, args...)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (g *Generator) generate() error {
	cTemplate := template.Must(template.New("nacos_template").Parse(configTemplate))

	var decodeFuncCode string
	if g.codec == "yaml" {
		decodeFuncCode = "err := gonacos.DecodeYaml(data, &c)"
	} else {
		decodeFuncCode = "err := gonacos.DecodeJson(data, &c)"
	}

	err := cTemplate.Execute(&g.buf, map[string]string{
		"ConfigName":     g.configName,
		"Group":          g.group,
		"DataId":         g.dataid,
		"DecodeFuncCode": decodeFuncCode,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Printf("format 错误: %s", err)
		return g.buf.Bytes()
	}
	return src
}
