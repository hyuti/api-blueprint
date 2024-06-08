package router

import "google.golang.org/protobuf/encoding/protojson"

var marshaller protojson.MarshalOptions

func init() {
	marshaller = protojson.MarshalOptions{EmitUnpopulated: true}
}
