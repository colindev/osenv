### env loader

load os env to struct

#### Example

```golang
package main

type Env struct {
    Int int `env:"int"`

    Omit map[string]interface{} `env:"-"`
}

func main(){
    var env Env
    if err := osenv.LoadTo(&env); err != nil {
        log.Fatal(err)
    }

    fmt.Println(env)
}

```
