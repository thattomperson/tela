package eddn

//go:generate schematyper -o blackmarket.go --package eddn ./schemas/blackmarket-v1.0.json

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/zeromq/goczmq"
)

type Schema struct {
	SchemaRef string `json:"$schemaRef"`
}

func Listen(endpoint string, handler func(interface{})) error {
	c, err := goczmq.NewSub(endpoint, "")
	if err != nil {
		return err
	}

	log.Println("Listening")

	for {
		bs, err := c.RecvMessage()
		if err != nil {
			return err
		}

		var buff []byte
		for _, _b := range bs {
			buff = append(buff, _b...)
		}

		b := bytes.NewReader(buff)
		r, err := zlib.NewReader(b)
		if err != nil {
			return err
		}

		buff, err = ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		var schema Schema

		json.Unmarshal(buff, &schema)

		switch schema.SchemaRef {
		case "https://eddn.edcd.io/schemas/commodity/3":
			var v Commodity
			json.Unmarshal(buff, &v)
			handler(v)
		case "https://eddn.edcd.io/schemas/journal/1":
		case "https://eddn.edcd.io/schemas/outfitting/2":
		case "https://eddn.edcd.io/schemas/shipyard/2":
		case "https://eddn.edcd.io/schemas/blackmarket/1":
		default:
			log.Printf("Unknown type %s\n", schema.SchemaRef)
		}

	}

	return nil
}
