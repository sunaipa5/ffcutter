package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/richtext"
	"github.com/ncruces/zenity"
)

var (
	vid_name_inp   nucular.TextEditor
	vid_folder_inp nucular.TextEditor

	cut_start  = "00:00:00"
	cut_end    = "00:00:00"
	cut_status int

	status   string
	no_audio bool
)

var rtxt *richtext.RichText

var time_ranges []string

var clean bool

func page_main(w *nucular.Window) {

	w.RowScaled(vh-47).Ratio(0.35, 0.65)
	if sw := w.GroupBegin("left", nucular.WindowBorder|nucular.WindowNoScrollbar); sw != nil {
		sw.RowScaled(20).Ratio(0.09, 0.09, 0.72, 0.1)
		for i, v := range time_ranges {
			var val bool

			if sw.SelectableLabel("U", "LC", &val) {
				if i > 0 {
					time_ranges[i], time_ranges[i-1] = time_ranges[i-1], time_ranges[i]
				}
			}

			if sw.SelectableLabel("D", "LC", &val) {
				if i < len(time_ranges)-1 {
					time_ranges[i], time_ranges[i+1] = time_ranges[i+1], time_ranges[i]
				}
			}

			sw.Label(v, "CC")

			if sw.SelectableLabel("-", "LC", &val) {
				time_ranges = append(time_ranges[:i], time_ranges[i+1:]...)
				i--
			}
		}

		sw.GroupEnd()
	}

	if sw := w.GroupBegin("right", nucular.WindowNoScrollbar); sw != nil {
		sw.RowScaled(w.WidgetBounds().H - 120).Dynamic(1)

		if swx := sw.GroupBegin("top", nucular.WindowNoScrollbar); swx != nil {

			swx.RowScaled(35).Dynamic(1)
			swx.Label("Add cut", "LC")

			swx.RowScaled(35).Ratio(0.425, 0.05, 0.425, 0.1)

			if sc := swx.Combo(label.T(cut_start), 250, nil); sc != nil {
				sc.Row(25).Dynamic(1)
				hour, minute, second, msecond := 0, 0, 0, 0

				parts := strings.Split(cut_start, ":")
				if len(parts) == 3 {
					fmt.Sscanf(parts[0], "%d", &hour)
					fmt.Sscanf(parts[1], "%d", &minute)
					secParts := strings.Split(parts[2], ".")
					fmt.Sscanf(secParts[0], "%d", &second)
					if len(secParts) > 1 {
						fmt.Sscanf(secParts[1], "%d", &msecond)
					}
				}

				sc.PropertyInt("Hour:", 0, &hour, 23, 1, 1)
				sc.PropertyInt("Minute:", 0, &minute, 59, 1, 1)
				sc.PropertyInt("Second:", 0, &second, 59, 1, 1)
				sc.PropertyInt("Millisecond:", 0, &msecond, 99, 1, 1)

				if msecond > 0 {
					cut_start = fmt.Sprintf("%02d:%02d:%02d.%02d", hour, minute, second, msecond)
				} else {
					cut_start = fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
				}
			}
			swx.Label("To", "LC")

			if sc := swx.Combo(label.T(cut_end), 250, nil); sc != nil {
				sc.Row(25).Dynamic(1)
				hour, minute, second, msecond := 0, 0, 0, 0

				parts := strings.Split(cut_end, ":")
				if len(parts) == 3 {
					fmt.Sscanf(parts[0], "%d", &hour)
					fmt.Sscanf(parts[1], "%d", &minute)
					secParts := strings.Split(parts[2], ".")
					fmt.Sscanf(secParts[0], "%d", &second)
					if len(secParts) > 1 {
						fmt.Sscanf(secParts[1], "%d", &msecond)
					}
				}

				sc.PropertyInt("Hour:", 0, &hour, 23, 1, 1)
				sc.PropertyInt("Minute:", 0, &minute, 59, 1, 1)
				sc.PropertyInt("Second:", 0, &second, 59, 1, 1)
				sc.PropertyInt("Millisecond:", 0, &msecond, 99, 1, 1)

				if msecond > 0 {
					cut_end = fmt.Sprintf("%02d:%02d:%02d.%02d", hour, minute, second, msecond)
				} else {
					cut_end = fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
				}
			}

			if swx.ButtonText("Add") {
				if !IsValidFFmpegTime(cut_start) {
					showPopup(w, "Error", "False format!")
					return
				}

				if !IsValidFFmpegTime(cut_end) {
					showPopup(w, "Error", "False format!")
					return
				} else {
					if cut_end == "00:00:00" || cut_end == "00:00:00.0" || cut_end == cut_start {
						showPopup(w, "Error", "You can't type same time!")
						return
					}
				}

				c := cut_start + "-" + cut_end
				if contains(time_ranges, c) {
					showPopup(w, "Warning", "This cut already exist")
					return
				}
				time_ranges = append(time_ranges, c)
			}

			swx.RowScaled(100).Dynamic(1)
			swx.Spacing(1)

			swx.GroupEnd()
		}

		if swx := sw.GroupBegin("bottom", nucular.WindowNoScrollbar); swx != nil {
			//Bottom

			swx.RowScaled(35).Dynamic(1)
			swx.CheckboxText("Disable Audio", &no_audio)

			swx.RowScaled(35).Ratio(0.15, 0.72, 0.13)
			swx.Label("Location:", "LC")
			vid_folder_inp.Edit(swx)
			vid_folder_inp.Flags = nucular.EditIbeamCursor
			vid_folder_inp.Active = false

			if swx.ButtonText("Change") {
				l, err := select_location()
				if err != nil {
					if err.Error() != "dialog canceled" {
						showPopup(w, "Error", "Failed to get location"+err.Error())
						return
					}
				}
				vid_folder_inp.Buffer = []rune(l)
			}

			swx.Label("Filename:", "LC")
			vid_name_inp.Flags = nucular.EditField | nucular.EditIbeamCursor
			vid_name_inp.Edit(swx)

			if swx.ButtonText("Cut") {
				if len(time_ranges) == 0 {
					status = "Please add a cut"
					return
				}

				if cut_status == 1 {
					return
				}

				go func() {
					cut_status = 1
					err := CutAndCombine(CutAndCombineOptions{
						InputFile:  video_location,
						OutputFile: filepath.Join(string(vid_folder_inp.Buffer), string(vid_name_inp.Buffer)),
						Ranges:     time_ranges,
						NoAudio:    no_audio,
						ProgressCB: func(p string) { status = p; w.Master().Changed() },
						ErrorCB:    func(msg string) { status = msg; w.Master().Changed() },
						FinishedCB: func() { status = "Done"; w.Master().Changed() },
					})

					if err != nil {
						showPopup(w, "Error", "Failed cut: "+err.Error())
					}
					cut_status = 0
				}()
			}

			swx.GroupEnd()
		}

		sw.GroupEnd()
	}
	w.Row(15).Dynamic(1)
	w.Label(status, "LC")
}

func select_location() (string, error) {
	f, err := zenity.SelectFile(
		zenity.Directory(),
		zenity.Title("Select Folder"),
	)

	if err != nil {
		return "", err
	}

	return f, nil
}
