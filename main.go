package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "bufio"
    "unicode"
    "github.com/gookit/color"
)

var (
    MAX_LINES_FUNCTION = 25
    MAX_FUNCTIONS = 10
    MAX_GLOBAL_FUNCTIONS = 5
    MAX_GLOBAL_VARIABLE = 0
)

type results struct {
    nb_functions_in_file int
    nb_functions_static int
    nb_functions_global int
    nb_global_variables int
}

func check_error(err error) {
    if (err != nil) {
        panic(err)
    }
}

/* Error: -1 | Directory: 0 | File: 1 */
func get_file_state(filename string) int8 {
    info, err := os.Stat(filename)

    if os.IsNotExist(err) {
        return -1
    } else if info.IsDir() {
        return 0;
    } else {
        return 1;
    }
}

func is_function_proto(line string) bool {
    i := 0
    for ; i < len(line); i++ {
        if line[i] == '(' {
            break
        }
    }

    if i >= len(line) || line[i] != '(' {
        return false
    }

    return line[len(line) - 1] == ')'
}

func is_open_bracket(line string) bool {
    return len(line) > 0 && line[0] == '{'
}

func is_close_bracket(line string) bool {
    return len(line) > 0 && line[0] == '}'
}

func is_static(line string) bool {
    return line[:6] == "static"
}

func is_variable(line string) bool {
    for _, char := range line {
        if char == '=' {
            return true
        }
    }
    return false
}

func is_blank_line(line string) bool {
    for _, char := range line {
        if unicode.IsLetter(char) || unicode.IsDigit(char) {
            return false
        }
    }
    return true
}

func is_single_line_comment(line string) bool {
    i := 2
    for ; i < len(line); i++ {
        if line[i-2:i] == "//" || line[i-2:i] == "/*" {
            break
        }
    }

    if len(line) == 0 {
        return false
    } else if line[i-2:i] == "//" {
        return true
    } else if line[i-2:i] == "/*" {
        return len(line) >= 2 && line[len(line)-2:len(line)] == "*/"
    } else {
        return false
    }
}

func is_begin_comment(line string) bool {
    for i:= 2; i <= len(line); i++ {
        if line[i-2:i] == "/*" {
            return true
        }
    }
    return false
}

func is_end_comment(line string) bool {
    return len(line) >= 2 && line[len(line) - 2:len(line)] == "*/"
}

func check_function(fileScanner *bufio.Scanner) (bool, int, string) {
    function_line_nb := 0
    function_name := ""

    scope_static := false

    line := fileScanner.Text()
    if is_function_proto(line) {
        scope_static = is_static(line)
        function_name = line

        for fileScanner.Scan() {
            if is_open_bracket(fileScanner.Text()) {
                break
            }
        }
        if is_open_bracket(fileScanner.Text()) {
            for fileScanner.Scan() {
                line_2 := fileScanner.Text()
                if is_close_bracket(line_2) {
                    break
                }
                if !is_single_line_comment(line_2) && is_begin_comment(line_2) {
                    for fileScanner.Scan() {
                        line_2 = fileScanner.Text()
                        if is_end_comment(line_2) {
                            break
                        }
                    }
                    fileScanner.Scan()
                }
                line_2 = fileScanner.Text()
                if !is_blank_line(line_2) && !is_single_line_comment(line_2) {
                    function_line_nb++
                }
            }
        }
        return scope_static, function_line_nb, function_name
    }

    return false, -1, ""
}

func print_results(file_results results, verbose bool) {
    if verbose {
        fmt.Println()
    }
    if file_results.nb_global_variables > MAX_GLOBAL_VARIABLE {
        color.Red.Printf("Variables --> Global: %d\n", file_results.nb_global_variables)
    }
    if verbose {
        fmt.Println()
    }
    if file_results.nb_functions_in_file > MAX_FUNCTIONS {
        color.Red.Printf("Functions -> Total: %d\n", file_results.nb_functions_in_file)
    } else if verbose {
        fmt.Printf("Functions -> Total: %d\n", file_results.nb_functions_in_file)
    }
    if file_results.nb_functions_global > MAX_GLOBAL_FUNCTIONS {
        color.Red.Printf("Functions -> Global: %d\n", file_results.nb_functions_global)
    } else if verbose {
        fmt.Printf("Functions -> Global: %d\n", file_results.nb_functions_global)
    }
    if verbose {
        fmt.Printf("Functions -> Static: %d\n", file_results.nb_functions_static)
        fmt.Println()
    }
}

func check_style_file(filename string, verbose bool) {
    file, err := os.Open(filename)
    defer file.Close()
    check_error(err)

    fmt.Printf("----- %s -----\n", filename)
    file_results := results{0, 0, 0, 0}

    fileScanner := bufio.NewScanner(file)

    for fileScanner.Scan() {
        scope_static, nb_line, func_name := check_function(fileScanner)
        if nb_line > 0 {
            if nb_line > MAX_LINES_FUNCTION {
                color.Red.Printf("%s: %d\n", func_name, nb_line)
            } else if verbose {
                fmt.Printf("%s: %d\n", func_name, nb_line)
            }
            file_results.nb_functions_in_file++
            if scope_static {
                file_results.nb_functions_static++
            } else {
                file_results.nb_functions_global++
            }
        } else {
            line := fileScanner.Text()
            if is_variable(line) && !is_static(line) {
                file_results.nb_global_variables++;
            }
        }
    }

    check_error(fileScanner.Err())
    print_results(file_results, verbose)
}

func handle_args(name string, verbose bool) {
    state := get_file_state(name)

    if state == -1 {
        return
    } else if state == 0 {
        files, err := ioutil.ReadDir(name)
        check_error(err)

        for _, file := range files {
            if name[len(name) - 1] == '/' {
                handle_args(name + file.Name(), verbose)
            } else {
                handle_args(name + "/" + file.Name(), verbose)
            }
        }
    } else {
        if name[len(name)-2:] == ".c" || name[len(name)-2:] == ".h" {
            check_style_file(name, verbose)
        }
    }
}

func main() {
    begin := 1
    verbose := false

    if os.Args[1] == "-v" || os.Args[1] == "--verbose" {
        begin++
        verbose = true
    }

    for _, arg := range os.Args[begin:] {
        handle_args(arg, verbose)
    }
}
