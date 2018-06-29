package loaders

import (
	"context"
	"fmt"
	"log"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"metagraf/internal/metagraf"
)

func CtoDgraph( mg metagraf.MetaGraf ) {
	dgraph := "127.0.0.1:9080"
	conn, err := grpc.Dial(dgraph, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()


	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))


	resp, err := dg.NewTxn().Query(context.Background(), `{
  bladerunner(func: eq(name@en, "Blade Runner")) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}`)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n", resp.Json)
}