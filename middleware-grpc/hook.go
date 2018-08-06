package main

import (
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

func MyRabbitHook(object *coprocess.Object) (*coprocess.Object, error) {
	return object, nil
}
