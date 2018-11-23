# go-dhallconfig

Currently a small go [dhall](https://github.com/dhall-lang/dhall-lang) config library.

Because your config is a dhall expression:

- Load config from environment variables.
- Load config from local files.
- Safely load config remote from a server.
- Factor your config with functions (while still collapsing to a known type.).
- Get type checking and nice error messages for your config.
- Get an automatic config formatter with ```dhall format```.

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
