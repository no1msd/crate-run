class_name Constants
extends Node

const LEVEL_COUNT := 50
const BASE_URL: String = "https://crate.run/"

var _dynamic_base_url: String = ""


func _ready() -> void:
	if OS.has_feature("JavaScript"):
		var result = JavaScript.eval("getEnvironmentVariables().basePath", true)
		if result:
			_dynamic_base_url = String(result)


func get_base_url() -> String:
	if _dynamic_base_url != "":
		return _dynamic_base_url
	
	return BASE_URL
