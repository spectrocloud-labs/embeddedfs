# embeddedfs
The embeddedfs package reads files and renders templates from an embedded file system.

## Usage
Assuming you have the following directory structure:

```console
.
├── main.go
└── myresources
    └── mytemplate.tmpl
```

And `mytemplate.tmpl` defined as follows:
```
Name: {{ .Name }}
```

Then the following code will produce `out.yaml` with the content `Name: Bob`:
```go
//go:embed resources/*
var resources embed.FS

func main() {
    e := embeddedfs.NewEmbeddedFS("myresources", resources)
    args := map[string]interface{}{
        "Name": "Bob",
    }
    if err := e.RenderTemplate(args, "", "mytemplate.tmpl", "out.yaml"); err != nil {
        return fmt.Errorf("failed to render my template: %w", err)
    }
}
```
