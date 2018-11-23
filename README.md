# go-dhallconfig

Currently a small go [dhal](https://github.com/dhall-lang/dhall-lang) config library.

Because your config is a dhall expression, you can safely load environment variables, fetch
config from remote sources, or factor your config with functions, but the end result is always
a statically typed value you can inspect.

It lets you type check and load dhall config files into go structs, if you config is of the wrong type, you
will get an error message.

# Example

```
type Config struct {
	Foo  string
	Bar  int64
	Baz  bool
	Boo  float64
	Bag  []int64
	Nest struct{ Foo bool }
}

...

	config := &Config{}

    configText := `
        { Foo =
            env:FOO as Text
        , Bar =
            +6
        , Baz =
            True
        , Boo =
            10.2
        , Bag =
            [ +1, +2 ]
        , Nest =
            { Foo = True }
        }
    `

	err := dhallconfig.LoadConfig(configText, actualConfig)
	if err != nil {
		panic(err)
	}
```

# Tips

You don't need to put any file loading logic in your code. Just accept a dhall expression from the command line.
"./foo.dhall" is a valid expression and the dhall interpreter will load your config.