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

	"golang.org/x/tools/go/packages"
)

var (
	structAtomicTemplate = `
var __{{ .ConfigName.StructFull }} = atomic.Value{}

func Register{{ .ConfigName.StructFull }}(ncc *gonacos.NacosConfigClient, defConf *{{ .ConfigName.StructFull }}) {
	if defConf != nil {
		__{{ .ConfigName.StructFull }}.Store(*defConf)
	}

	ncc.WatchF("{{ .Group }}", "{{ .DataId }}", func(namespace, group, dataId, data string) {
		var c {{ .ConfigName.StructFull }}
		{{ .DecodeFuncCode }}
		if err == nil {
			__{{ .ConfigName.StructFull }}.Store(c)
		}
	})
}

func Get{{ .ConfigName.StructFull }}() {{ .ConfigName.StructFull }} {
	return __{{ .ConfigName.StructFull }}.Load().({{ .ConfigName.StructFull }})
}
`

	structEmbedTemplate = `
var {{ .MutexVar }} = sync.RWMutex{}

func RegisterEmbed{{ .EmbedStruct }}{{ .ConfigName.VarName }}(ncc *gonacos.NacosConfigClient, c *{{ .ConfigName.StructFull }}) {
	ncc.WatchF("{{ .Group }}", "{{ .DataId }}", func(namespace, group, dataId, data string) {
		{{ .MutexVar }}.Lock()
		defer {{ .MutexVar }}.Unlock()

		{{ .DecodeFuncCode }}
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}

func (c *{{ .EmbedStruct }}) Get{{ .ConfigName.VarName }}() {{ .ConfigName.StructFull }} {
	{{ .MutexVar }}.RLock()
	defer {{ .MutexVar }}.RUnlock()

	return *c.{{ .ConfigName.VarName }}
}
`
)

var (
	// 配置类型，空则是标准类型，带参数，这是该参数的嵌入类型
	embedStruct = flag.String("embed", "", "嵌入式配置的主结构体名称")

	// 配置信息
	configName = flag.String("config", "", "结构体名称")
	dataid     = flag.String("dataid", "", "nacos内的dataid名称")
	group      = flag.String("group", "", "nacos内的group名称，不填则使用NacosConfigClient的defaultGroup")
	codec      = flag.String("codec", "yaml", "编码格式，支持json、yaml")
)

type FlagConfig struct {
	VarName     string
	StructFull  string
	PackageName string
	StructName  string
}

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
		embedStruct: *embedStruct,
		configName:  getConfigName(*configName),
		dataid:      *dataid,
		group:       *group,
		codec:       *codec,
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"gonacos_config %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("package %s", currentPackage)
	g.Printf("\n")

	// 分析
	g.parseImport()

	// 生成代码
	err = g.generate()
	if err != nil {
		log.Println("生成失败：" + err.Error())
		return
	}

	// 格式化并生成代码
	src := g.format()

	// Write to file.
	var baseName string
	if embedStruct == nil || *embedStruct == "" {
		baseName = fmt.Sprintf("nacos_%s.gen.go", strings.ToLower(*dataid))
	} else {
		baseName = fmt.Sprintf("embed_%s.gen.go", strings.ToLower(*dataid))
	}
	outputName := filepath.Join(dir, strings.ToLower(baseName))
	err = os.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("写文件失败: %s", err)
	}
}

func getConfigName(configName string) FlagConfig {
	split := strings.Split(configName, ":")

	if len(split) == 1 {
		p, s := getStructInfo(split[0])
		return FlagConfig{
			VarName:     split[0],
			StructFull:  split[0],
			PackageName: p,
			StructName:  s,
		}
	} else {
		p, s := getStructInfo(split[1])
		return FlagConfig{
			VarName:     split[0],
			StructFull:  split[1],
			PackageName: p,
			StructName:  s,
		}
	}
}

func getStructInfo(structName string) (string, string) {
	ss := strings.Split(structName, ".")
	if len(ss) == 1 {
		return "", structName
	} else {
		return ss[0], ss[1]
	}
}

type Generator struct {
	buf bytes.Buffer

	embedStruct string

	configName FlagConfig
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

func (g *Generator) parseImport() {
	// 固定包
	g.Printf("import gonacos \"github.com/reatang/go-simple-nacos\"\n")

	// 分析其他包的配置
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedImports,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, "./"+os.Getenv("GOFILE"))

	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, p := range pkg.Imports {
			if p.Name == g.configName.PackageName {
				g.Printf("import \"%s\"\n", p.PkgPath)
			}
		}
	}
}

func (g *Generator) generate() error {
	var cTemplate *template.Template
	if g.embedStruct != "" {
		g.Printf("import \"fmt\"\n")
		g.Printf("import \"sync\"\n")
		cTemplate = template.Must(template.New("embed_template").Parse(structEmbedTemplate))
	} else {
		g.Printf("import \"sync/atomic\"\n")
		cTemplate = template.Must(template.New("nacos_template").Parse(structAtomicTemplate))
	}

	var decodeFuncCode string
	if g.codec == "yaml" {
		decodeFuncCode = "err := gonacos.DecodeYaml(data, &c)"
	} else {
		decodeFuncCode = "err := gonacos.DecodeJson(data, &c)"
	}

	err := cTemplate.Execute(&g.buf, map[string]interface{}{
		"EmbedStruct":    g.embedStruct,
		"ConfigName":     g.configName,
		"Group":          g.group,
		"DataId":         g.dataid,
		"DecodeFuncCode": decodeFuncCode,
		"MutexVar":       fmt.Sprintf("__%s%sMutex", g.embedStruct, g.configName.VarName),
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
