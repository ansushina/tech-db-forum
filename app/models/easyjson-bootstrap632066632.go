// +build ignore

// TEMPORARY AUTOGENERATED FILE: easyjson bootstapping code to launch
// the actual generator.

package main

import (
  "fmt"
  "os"

  "github.com/mailru/easyjson/gen"

  pkg "github.com/ansushina/tech-db-forum/app/models"
)

func main() {
  g := gen.NewGenerator("model_posts_easyjson.go")
  g.SetPkg("models", "github.com/ansushina/tech-db-forum/app/models")
  g.Add(pkg.EasyJSON_exporter_Posts(nil))
  if err := g.Run(os.Stdout); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
