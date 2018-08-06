package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/MinecraftXwinP/diplomat"
	"github.com/google/subcommands"
)

type createCommand struct {
	out string
}

func (*createCommand) Name() string {
	return "create"
}

func (*createCommand) Synopsis() string {
	return "generate example outline file"
}

func (*createCommand) Usage() string {
	return `create: generate example outline file`
}
func (c *createCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.out, "out", "outline.yaml", "filepath of generated outline file")
}

const exampleOutlineFile = `version: '1'
settings:
  chinese:
    convert:
      mode: t2s
      from: zh-TW
      to: zh-CN
  copy:
  - from: en
    to: fr
fragments:
  admin:
    description: translations for admin page
    translations:
      admin:
        zh-TW: 管理員
        en: Admin
output:
  fragments:
  - type: js
    name: "{{.Locale}}.{{.FragmentName}}.js"`

func (c *createCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("cannot get current working directory", err)
		return subcommands.ExitFailure
	}
	path := filepath.Join(wd, c.out)
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 755)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	_, err = os.Stat(path)
	if err == nil {
		log.Printf("%s already exist.", path)
		return subcommands.ExitFailure
	}
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
		return subcommands.ExitFailure
	}
	outF, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	defer outF.Close()
	_, err = outF.WriteString(exampleOutlineFile)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

type generateCommand struct {
	outlineFile string
	outdir      string
	watch       bool
}

func (*generateCommand) Name() string {
	return "generate"
}

func (*generateCommand) Synopsis() string {
	return "generate language modules"
}

func (*generateCommand) Usage() string {
	return "generate: generate language modules according to outline file"
}

func (g *generateCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&g.outlineFile, "f", "outline.yaml", "path to outline file")
	f.StringVar(&g.outdir, "out", "out", "output dir")
	f.BoolVar(&g.watch, "watch", false, "watch file changes")
}

func (g *generateCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	initFunc := diplomat.NewDiplomatForFile
	if g.watch {
		initFunc = diplomat.NewDiplomatWatchFile
	}
	d, err := initFunc(g.outlineFile, g.outdir)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	d.RegisterMessenger("js", diplomat.JsModuleMessengerHandler)
	if g.watch {
		d.Watch()
	} else {
		err = d.Output()
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	}

	return subcommands.ExitSuccess
}