package main

import (
	"fmt"
	"os"
)

func main () {
	file, _ := os.Create("test.go")
	codeSlice := []string{
		"package main\n",
		"import (\n",
		"\t",
		`"fmt"`,
		"\n)\n",
		"func main () {\n\t",
		` fmt.Println("asd")`,
		"router := gin.Default()",
		"",
		"\n}",
	}
	for _, v := range codeSlice {
		_, _ = fmt.Fprintf(file,"%s", v)

	}
	_ = file.Close()
}
