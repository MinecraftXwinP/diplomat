package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
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

const exampleOutlineFile = `version: "1"
preprocessors:
- type: chinese
  options:
    - mode: t2s
      from: zh-TW
      to: zh-CN
output:
  - selectors:
      - admin
      - manage
    templates:
      - type: js
        options:
          filename: "{{.Lang}}.locale.js"
`

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
	folder string
	outdir string
	watch  bool
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
	f.StringVar(&g.folder, "dir", "diplomat", "path to diplomat folder")
	f.StringVar(&g.outdir, "out", "out", "output dir")
	f.BoolVar(&g.watch, "watch", false, "watch file changes")
}
func (g *generateCommand) doWatch(c context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	d, errChan, changeListener := diplomat.NewDiplomatWatchDirectory(g.folder)
	go func() {
		for e := range errChan {
			log.Println("error:", e)
		}
	}()
	go func() {
		for range changeListener {
			d.Output(g.outdir)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return subcommands.ExitSuccess
}

func (g *generateCommand) Execute(c context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if g.watch {
		return g.doWatch(c, f)
	}
	d, err := diplomat.NewDiplomatForDirectory(g.folder)
	if err != nil {
		log.Println("error:", err)
		return subcommands.ExitFailure
	}
	err = d.Output(g.outdir)
	if err != nil {
		log.Println("error:", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}