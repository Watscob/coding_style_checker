# Coding Style Checker
A tool to check certain features of the coding style in C files.

**Warning**
If you put a file which is not a .c or .h file, the project will ignore it.
**Warning**
Do not use regex and expandable syntax such as `my_dir/*` or `*.c`

## Dependencies
- [ ] [go](https://golang.org/dl/)

## Instructions to use the project without compiling it
### Usage
Go in the cloned directory and execute the following command
```sh
go run main.go <directories/filenames>

# <directories/filenames> is a list of directories and/or filenames separed by one or more spaces

# Example:
go run main.go test.c

# or
go run main.go my_dir
```

## Instructions to use the compiled project
### Build
Go in the cloned directory and execute the following command
```sh
go build main.go
```

### Usage
```sh
./main.go <directories/filenames>

# <directories/filenames> is a list of directories and/or filenames separed by one or more spaces

# Example:
./main.go test.c

# or
./main.go my_dir
```
