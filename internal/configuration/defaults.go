package configuration

func defaults() map[Key]interface{} {
	defaultValues := make(map[Key]interface{})
	defaultValues[KeyPort] = 8080
	return defaultValues
}
