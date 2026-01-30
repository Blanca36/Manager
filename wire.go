//go:build wireinject
// +build wireinject

package main

//wire依赖注入
import (
	"github.com/google/wire"
)

func InitializeEvent() Event {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}
}
