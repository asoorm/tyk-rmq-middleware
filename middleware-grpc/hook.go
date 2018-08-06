package main

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

func MyRabbitHook(object *coprocess.Object) (*coprocess.Object, error) {

	res, err := doRPC(object.Request.Body)
	if err != nil {

		type ErrorFormat struct {
			Message string
			Error   string
		}

		errorJson, _ := json.Marshal(ErrorFormat{
			Message: "failure doing rabbit rpc",
			Error:   err.Error(),
		})

		object.Request.ReturnOverrides.ResponseCode = 666
		object.Request.ReturnOverrides.ResponseError = string(errorJson)
		logrus.WithError(err).Error("failure doing rabbit rpc")

		return object, nil
	}

	object.Request.ReturnOverrides.ResponseCode = 200
	object.Request.ReturnOverrides.ResponseError = res

	logrus.Info("success!")

	return object, nil
}
