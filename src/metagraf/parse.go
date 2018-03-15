package metagraf

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func Parse( filepath string ) {
	b, err := ioutil.ReadFile( filepath )
	if err != nil {
		panic(err)
	}

	var mg MetaGraf

	json.Unmarshal(b, &mg)
	if err != nil {
		panic(err)
	}

	fmt.Printf( "%T", mg )
	fmt.Println( mg )

}
