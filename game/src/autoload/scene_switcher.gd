class_name SceneSwitcher
extends Node

var _next_scene: String = ""
var _params = null
var _timer: Timer = Timer.new()


func _ready() -> void:
	add_child(_timer)
	_timer.connect("timeout", self, "on_timer_ended")


func on_timer_ended() -> void:
	if _next_scene != "":
		get_tree().change_scene(_next_scene)

	_next_scene = ""
	_timer.stop()


func change_scene(next_scene: String, timeout: float, params = null) -> void:
	get_tree().get_current_scene().animate_out()
	_next_scene = next_scene
	_params = params
	_timer.start(timeout)


func get_param(name):
	if _params != null and _params.has(name):
		return _params[name]

	return null
