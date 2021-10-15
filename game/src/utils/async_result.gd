class_name AsyncResult
extends Node

signal finished(async_result)

var _finished := false
var _result = null
var _error = null
var _failing := false


func is_failing() -> bool:
	return _failing


func is_finished() -> bool:
	return _finished


func get_result():
	return _result


func get_error():
	return _error


func is_error() -> bool:
	return _finished and _error != null


func success(result) -> void:
	_result = result
	_finished = true
	emit_signal("finished", self)


func failure(error) -> void:
	_error = error
	_failing = true
	_finished = true
	emit_signal("finished", self)


func deferred_failure(error) -> void:
	_failing = true
	call_deferred("failure", error)
