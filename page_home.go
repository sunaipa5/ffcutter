package main

import (
	"path/filepath"

	"github.com/aarzilli/nucular"
	"github.com/ncruces/zenity"
)

var video_location string

func page_home(w *nucular.Window) {
	vh := w.LayoutAvailableHeight()

	w.Row(vh / 3).Dynamic(1)
	w.Row(35).Ratio(0.2, 0.6, 0.2)
	w.Spacing(1)
	if w.ButtonText("Select a video") {
		fp := selectVideo(w)
		if fp != "" {
			//video_location = "/home/user/Videos/example.mp4"
			video_location = fp

			base := filepath.Base(fp)
			ext := filepath.Ext(fp)
			name := base[:len(base)-len(ext)]

			//vid_folder_inp.Buffer = []rune("/home/user/Videos/")
			vid_folder_inp.Buffer = []rune(filepath.Dir(fp))
			vid_name_inp.Buffer = []rune(name + "(trimmed)" + ext)
			currentPage = 2
		}
	}
	w.Spacing(1)

}

func selectVideo(w *nucular.Window) string {
	filepath, err := zenity.SelectFile(
		zenity.Title("Select Video File"),
		zenity.FileFilters{
			{
				Name:     "Video File",
				Patterns: []string{"*.mp4", "*.mkv", "*.avi", "*.mov", "*.webm"},
				CaseFold: true,
			},
		},
	)

	if err != nil {
		if err.Error() != "dialog canceled" {
			showPopup(w, "Error", "Failed to get video file "+err.Error())
			return ""
		}
	}

	return filepath
}
