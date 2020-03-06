package output

import (
	"log"
	"os"
	"os/exec"

	"github.com/ma6254/FictionDown/store"
)

type PandocEPUB struct {
}

func (t *PandocEPUB) Conv(src store.Store, outpath string, opts Option) (err error) {

	mdpath := outpath + ".md"

	if err = Output(src, "md", mdpath, opts); err != nil {
		return
	}

	log.Printf("中间文件转换完成: %#v\n", mdpath)
	c := exec.Command(
		"pandoc",
		"--epub-chapter-level", "2",
		"-f", "markdown-raw_tex",
		"-o", outpath,
		mdpath)
	log.Printf("调用Pandoc: %#v %#v\n", c.Path, c.Args)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}
	return os.Remove(mdpath)
}
