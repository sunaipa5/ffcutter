package main

import (
	"path/filepath"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/style"
)

var activeTheme style.Theme

func menubar(w *nucular.Window) {
	w.MenubarBegin()

	w.Row(20).Ratio(0.2, 0.8)
	if w := w.Menu(label.TA("FFcutter", "LC"), 120, nil); w != nil {
		w.Row(25).Dynamic(1)
		if w.MenuItem(label.TA("Select a Video", "LC")) {
			fp := selectVideo(w)
			if fp != "" {
				video_location = fp

				base := filepath.Base(fp)
				ext := filepath.Ext(fp)
				name := base[:len(base)-len(ext)]

				vid_folder_inp.Buffer = []rune(filepath.Dir(fp))
				vid_name_inp.Buffer = []rune(name + "(trimmed)" + ext)
				currentPage = 2
			}
		}
		if w.MenuItem(label.TA("About", "LC")) {
			showPopup(w, "Info", "FFcutter v0.1 alpha")
		}
	}

	w.Label(video_location, "RC")

}
