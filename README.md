### env loader

load os env to struct

#### Example

```golang
package main

type Env struct {
    Path string `env:"PATH"`
    User string `env:"USER"`
    DefaultValue bool `env:"DV,true"`
    CustomInt int `env:"custom_int"`
    Omit map[string]interface{} `env:"-"`
}

func init(){
    os.Setenv("custom_int", "123")
}

func main(){
    var env Env
    if err := osenv.LoadTo(&env); err != nil {
        log.Fatal(err)
    }

    fmt.Println(env)
}
```
