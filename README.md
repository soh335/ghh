# ghh

ghh is cui tool for [github repository hook](http://developer.github.com/v3/repos/hooks/).

ghh use ```EDITOR``` environment variable for creating/editing hook.

# how to install

```
go install github.com/soh335/ghh
```

# usage

## set token

```
ghh config token <token>
```

you can get <token> at https://github.com/settings/applications .
config file is saved at  ```$HOME/.config/ghh```.

## support type

```
ghh support
```

## list hooks

```
ghh list <owner> <repo>
```

## show hook

```
ghh show <owner> <repo> <id>
```

## create hook

```
ghh create <owner> <repo> <type>
```

## edit hook

```
ghh edit <owner> <repo> <id>
```

## delete hook

```
ghh delete <owner> <repo> <id>
```

## test hook

```
ghh test <owner> <repo> <id>
```
