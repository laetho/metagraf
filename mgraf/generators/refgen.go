package metagraf

import (
	"fmt"
	"metagraf/mgraf/metagraf"
)

// Takes a pointer to a metagraf.MetaGraf struct
func Refgen( mg metagraf.MetaGraf ) {
	fmt.Println( mg )
	fmt.Println( mg.Kind, mg.Metadata.Name, mg.Metadata.Version)
}
