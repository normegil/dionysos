package model

type Resource string

const (
	ResourceItem    = Resource("item")
	ResourceUser    = Resource("user")
	ResourceStorage = Resource("storage")
)

type Action string

const (
	ActionRead  = Action("read")
	ActionWrite = Action("write")
)
