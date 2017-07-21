### env loader

load os env to struct

[![PayPal donate button](https://img.shields.io/badge/paypal-donate-yellow.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&hosted_button_id=YL3P2TXPF2GKE&item_name=node%2ddota2%2dspectator&currency_code=USD)
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
