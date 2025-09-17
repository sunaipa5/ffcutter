package main

import (
	_ "embed"
	"image"
	"image/color"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/font"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/aarzilli/nucular/style"
)

//go:embed Cantarell-Regular.ttf
var cantarell []byte

var (
	currentPage int
	vw          int
	vh          int
)

func main() {

	wnd := nucular.NewMasterWindowSize(0, "FFcutter", image.Point{700, 400}, mainLoop)
	wnd.SetStyle(style.FromTheme(style.DefaultTheme, 1.2))

	wnd.SetStyle(style.FromTable(style.ColorTable{
		ColorText:         color.RGBA{200, 200, 200, 255},
		ColorWindow:       color.RGBA{36, 36, 36, 255},
		ColorButton:       color.RGBA{0, 99, 0, 255},
		ColorButtonHover:  color.RGBA{0, 120, 0, 255},
		ColorButtonActive: color.RGBA{0, 120, 0, 255},
		ColorToggle:       color.RGBA{0, 99, 0, 255},
		ColorToggleHover:  color.RGBA{0, 120, 0, 255},
		ColorToggleCursor: color.RGBA{36, 36, 36, 255},
		ColorCombo:        color.RGBA{54, 54, 54, 255},
		ColorEdit:         color.RGBA{54, 54, 54, 255},
		ColorEditCursor:   color.RGBA{36, 36, 36, 255},
		ColorBorder:       color.RGBA{80, 80, 80, 255},
	}, 1.2))

	theming(wnd.Style())

	setFont(wnd)

	wnd.Main()
}

func setFont(wnd nucular.MasterWindow) {
	f, err := font.NewFace(cantarell, 16)
	if err != nil {
		panic(err)
	}
	wnd.Style().Font = f
}

func theming(s *style.Style) {
	s.NormalWindow.FooterPadding = image.Point{-15, -15}
	s.MenuWindow.FooterPadding = image.Point{0, 0}
	s.GroupWindow.FooterPadding = image.Point{0, 0}

	s.NormalWindow.Padding = image.Point{5, 5}
	s.MenuWindow.Padding = image.Point{0, 0}
	s.GroupWindow.Padding = image.Point{0, 0}

	s.Edit.Rounding = 7
	s.Combo.Rounding = 7
	s.Button.Border = 0
	s.Button.Rounding = 7
}

func mainLoop(w *nucular.Window) {
	vh = w.LayoutAvailableHeight()
	vw = w.LayoutAvailableWidth()

	switch currentPage {
	case 1:
		page_home(w)
	case 2:
		menubar(w)
		page_main(w)
	default:
		page_home(w)
	}
}

func showPopup(w *nucular.Window, popupHeader, popupLable string) {
	w.Master().PopupOpen(
		popupHeader,
		nucular.WindowTitle|nucular.WindowDynamic|nucular.WindowNoScrollbar,
		rect.Rect{
			X: vw/2 - 180,
			Y: vh/2 - 120,
			W: 250,
			H: 150,
		},
		true,
		func(w *nucular.Window) {
			w.Row(20).Dynamic(1)
			w.Label(popupLable, "LC")
			if w.Button(label.T("OK"), false) {
				w.Close()
			}
		},
	)
}
