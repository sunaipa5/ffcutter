## Shiny wm_class support for taskbar icon

> driver/x11driver/screen.go:460

```go
    class := "org.sunaipa.ffcutter"
    xproto.ChangeProperty(
    	s.xc,
    	xproto.PropModeReplace,
    	xw,
    	xproto.AtomWmClass,
    	xproto.AtomString,
    	8,
    	uint32(len(class)),
    	[]byte(class),
    )
```
