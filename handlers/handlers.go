package handlers

import (
	"github.com/ethanmidgley/the-sequel/in-memory/pkg/resp"
	"github.com/ethanmidgley/the-sequel/in-memory/store"
)

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG!"}
	}
	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	store.Store.Set(key, value)

	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	value := store.Store.Get(key)
	if value == "" {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}
