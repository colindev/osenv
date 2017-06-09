### env loader

load os env to struct

[![Build Status](https://travis-ci.org/colindev/osenv.svg?branch=master)](https://travis-ci.org/colindev/osenv)
[![GoDoc](https://godoc.org/github.com/colindev/osenv?status.svg)](https://godoc.org/github.com/colindev/osenv)

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
