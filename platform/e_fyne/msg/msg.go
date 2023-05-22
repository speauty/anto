package msg

import (
	"fmt"
	"github.com/sqweek/dialog"
)

func f() {
	dialog.Message("%s", "Please select a file").Title("Hello world!").Info()
	file, err := dialog.File().Title("Save As").Filter("All Files", "*").Save()
	fmt.Println(file)
	fmt.Println("Error:", err)
	dialog.Message("You chose file: %s", file).Title("Goodbye world!").Error()
	dialog.Directory().Title("Now find a dir").Browse()
}
