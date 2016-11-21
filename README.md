### env loader

load os env to struct

#### Example

```golang

type Env struct {
    Int int `env:"int"`

    Omit map[string]interface{} `env:"-"`
}

func main(){
    var ev Env
    if err := env.LoadTo(&ev); err != nil {
        log.Fatal(err)
    }

    fmt.Println(ev)
}

```
