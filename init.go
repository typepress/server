package server

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/achun/tom-toml"

	"github.com/typepress/core"
)

const usageMessage = `` +
	`Usage of TypePress:
Usage:
    typepress [-Falgs [Value]]

Flags:

    domain    ""                  hosting domain.
    laddr     ":80"               HTTP listen local network address.
    static    "static/"           path to static files.
    content   "content/"          path to content files.
    template  "template/"         path to template files. 
    config    ""                  path to config file.
    varsion                       prints Server module version.

Config file fist.
If config file is empty, try to load "conf/default.toml".
`

var (
	domain, laddr, contentPath, templatePath,
	staticPath, configFile string
	printsVarsion bool
)

func init() {
	flag.StringVar(&domain, "domain", "", "hosting domain")
	flag.StringVar(&laddr, "laddr", ":80", "HTTP listen local network address")
	flag.StringVar(&contentPath, "content", "content/", "path to content files")
	flag.StringVar(&templatePath, "template", "template/", "path to template files")
	flag.StringVar(&staticPath, "static", "static/", "path to static files")
	flag.StringVar(&configFile, "config", "", "path to config file")
	flag.BoolVar(&printsVarsion, "version", false, "prints current version")
}

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	os.Exit(2)
}

// +dl en
/*
  LoadConfig to load base configure from:
	- Command-line flags.
	- the environment variable, to sets such as: os.Setenv(core.SessionName+"_laddr").
	- TOML format config file.
*/
// +dl

/*
  LoadConfig 调入基本配置从:
	- 命令行参数.
	- 环境变量, 通过 os.Setenv(core.SessionName+"_laddr") 形式进行设置.
	- TOML 格式配置文件.
*/
func LoadConfig() {
	const try = "conf/default.toml"
	var err error
	flag.Usage = usage
	flag.Parse()
	if printsVarsion {
		fmt.Println(Version)
		os.Exit(0)
	}
	getEnv()

	config := try
	if len(configFile) != 0 {
		config = configFile
	}

	if !path.IsAbs(config) {
		config = path.Join(core.PWD, config)
	}

	core.Conf, err = toml.LoadFile(config)

	if err != nil && len(configFile) != 0 {
		fmt.Println("init site config:", err)
		os.Exit(1)
	}
	if err != nil {
		config = ""
		core.Conf = toml.Toml{}
	}

	conf := core.Conf
	if err != nil || nil == conf["site"] {
		conf["site"] = toml.NewItem(toml.Table)
	}
	if conf["template"] == nil {
		conf["template"] = toml.NewItem(toml.Table)
	}
	if conf["site.domain"] == nil {
		conf["site.domain"], _ = newItem(toml.String, domain)
	}
	if conf["site.laddr"] == nil {
		conf["site.laddr"], _ = newItem(toml.String, laddr)
	}
	if conf["site.content"] == nil {
		conf["site.content"], _ = newItem(toml.String, contentPath)
	}
	if conf["site.static"] == nil {
		conf["site.static"], _ = newItem(toml.String, staticPath)
	}
	if conf["template.path"] == nil {
		conf["template.path"], _ = newItem(toml.String, templatePath)
	}

	domain = conf["site.domain"].String()
	laddr = conf["site.laddr"].String()
	contentPath = conf["site.content"].String()
	staticPath = conf["site.static"].String()
	templatePath = conf["template.path"].String()

	fmt.Printf(`TypePress %s starting... server module %s

  listen   : %s
  domain   : %s
  static   : %s
  content  : %s
  template : %s
  config   : %s
`,
		AppVersion, Version, laddr, domain, staticPath, contentPath, templatePath, config)
}

func newItem(kind toml.Kind, x interface{}) (*toml.Item, error) {
	it := toml.NewItem(kind)
	return it, it.Set(x)
}

// Get flags from the environment variable
func getEnv() {
	var tmp string
	tmp = os.Getenv(core.SessionName + "_domain")
	if len(tmp) != 0 {
		domain = tmp
	}
	tmp = os.Getenv(core.SessionName + "_laddr")
	if len(tmp) != 0 {
		laddr = tmp
	}
	tmp = os.Getenv(core.SessionName + "_config")
	if len(tmp) != 0 {
		configFile = tmp
	}
	tmp = os.Getenv(core.SessionName + "_static")
	if len(tmp) != 0 {
		staticPath = tmp
	}
	tmp = os.Getenv(core.SessionName + "_content")
	if len(tmp) != 0 {
		contentPath = tmp
	}
	tmp = os.Getenv(core.SessionName + "_template")
	if len(tmp) != 0 {
		templatePath = tmp
	}
}
