class_name LoadingManager
extends Node

signal loading_changed(loading)

var _loading_count: int = 0
var _mtx: Mutex = Mutex.new()


func is_loading() -> bool:
	_mtx.lock()
	var is_loading := _loading_count > 0
	_mtx.unlock()

	return is_loading


func start_loading() -> void:
	_mtx.lock()

	_loading_count += 1

	if _loading_count == 1:
		emit_signal("loading_changed", true)

	_mtx.unlock()


func stop_loading() -> void:
	_mtx.lock()

	_loading_count -= 1

	if _loading_count == 0:
		emit_signal("loading_changed", false)

	_mtx.unlock()
