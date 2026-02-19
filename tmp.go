package main

import (
    "fmt"
    "goMud/internal/gmsl/lexer"
    "goMud/internal/gmsl/parser"
    "goMud/internal/gmsl/compiler"
    "os"
)

func main() {
    b, err := os.ReadFile("mudlib/locations/room_a.gms")
    if err != nil {
        panic(err)
    }
    
    l := lexer.NewLexer(string(b))
    p := parser.NewParser(l)
    ast := p.Parse() // Ignore the string it prints
    
    fmt.Println("\n=> DEBUG AST:", len(ast.Functions), "functions")
    if len(ast.Functions) > 0 {
        fmt.Println("=> DEBUG FUNCS:", ast.Functions[0].Name.Value, "Statements:", len(ast.Functions[0].Statements))
    }
    
    aout := compiler.NewCompiler(ast).Compile()
    
    fmt.Println("\n--- ASSEMBLY OUTPUT ---")
    fmt.Println(aout.String())
}
    fmt.Println("Imports len:", len(ast.Imports))
package main

import (
    "fmt"
    "goMud/internal/gmsl/lexer"
    "goMud/internal/gmsl/parser"
    "goMud/internal/gmsl/compiler"
    "os"
)

func main() {
    b, err := os.ReadFile("mudlib/locations/room_a.gms")
    if err != nil {
        panic(err)
    }
    
    l := lexer.NewLexer(string(b))
    p := parser.NewParser(l)
    ast := p.Parse() // Ignore the string it prints
    
    fmt.Println("\n=> DEBUG AST:", len(ast.Functions), "functions")
    if len(ast.Functions) > 0 {
        fmt.Println("=> DEBUG FUNCS:", ast.Functions[0].Name.Value, "Statements:", len(ast.Functions[0].Statements))
    }
    
    aout := compiler.NewCompiler(ast).Compile()
    
    fmt.Println("\n--- ASSEMBLY OUTPUT ---")
    fmt.Println(aout.String())
}
