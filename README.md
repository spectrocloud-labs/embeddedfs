[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/spectrocloud-labs/embeddedfs/issues)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Test](https://github.com/spectrocloud-labs/embeddedfs/actions/workflows/ci.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/spectrocloud-labs/embeddedfs)](https://goreportcard.com/report/github.com/spectrocloud-labs/embeddedfs)
[![codecov](https://codecov.io/gh/spectrocloud-labs/embeddedfs/graph/badge.svg?token=QZXWTOZ1LM)](https://codecov.io/gh/spectrocloud-labs/embeddedfs)
[![Go Reference](https://pkg.go.dev/badge/github.com/spectrocloud-labs/embeddedfs.svg)](https://pkg.go.dev/github.com/spectrocloud-labs/embeddedfs)

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
