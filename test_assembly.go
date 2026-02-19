package main

import (
    "fmt"
    "goMud/internal/gmsl/lexer"
    "goMud/internal/gmsl/parser"
    "goMud/internal/gmsl/compiler"
    "os"
)

func main() {
    b, err := os.ReadFile("mudlib/player_handler.gms")
    if err != nil {
        panic(err)
    }
    
    l := lexer.NewLexer(string(b))
    p := parser.NewParser(l)
    ast := p.Parse() // Ignore the string it prints
    
    aout := compiler.NewCompiler(ast).Compile()
    
    fmt.Println("--- ASSEMBLY OUTPUT ---")
    fmt.Println(aout.String())
}
