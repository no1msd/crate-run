class_name LevelCompleteDialog
extends ColorRect

signal close(dialog, holder)

var _holder: CanvasLayer = null
var _level_number := 0
var _moves := 0
var _timer := 0
var _prev_best_moves := 0
var _prev_best_timer := 0

var _callback_obj
var _callback_func


func _ready():
	anim_helper.fade_in(self)
	anim_helper.slide_in_from_down($Control)


func set_holder(holder: CanvasLayer) -> void:
	_holder = holder


func set_callback(callback_obj, callback_func) -> void:
	_callback_obj = callback_obj
	_callback_func = callback_func


func set_arguments(arguments) -> void:
	_level_number = arguments[0]
	_moves = arguments[1]
	_timer = arguments[2]
	var formatted_time = arguments[3]
	_prev_best_moves = arguments[4]
	_prev_best_timer = arguments[5]
	var prev_formatted_time = arguments[6]

	var completed = $Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/VBoxContainer/CompletedLabel
	var stats = $Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/VBoxContainer/StatisticsLabel
	var best = $Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/VBoxContainer/PersonalBestLabel

	completed.text = "Level " + str(_level_number) + " completed!"
	stats.text = "Moves: " + str(_moves) + "\nTime: " + formatted_time

	if not auth_context.is_logged_in():
		best.visible = false
	
	if _prev_best_moves > _moves or (_prev_best_moves == _moves and _prev_best_timer > _timer):
		best.text = "New personal best!"
	elif _prev_best_moves > 0:
		best.text = (
			"Previous personal best: "
			+ str(_prev_best_moves)
			+ " moves, "
			+ prev_formatted_time
		)


func _on_close() -> void:
	anim_helper.fade_out(self)
	anim_helper.slide_out_down($Control)

	yield(get_tree().create_timer(0.2), "timeout")

	emit_signal("close", self, _holder)


func _on_SelectorButton_pressed() -> void:
	if auth_context.is_logged_in():
		yield(highscore_api_service.register_highscore(_level_number, _moves, _timer), "finished")
		
	funcref(_callback_obj, _callback_func).call_func(false)
	_on_close()


func _on_NextLevelButton_pressed() -> void:
	if auth_context.is_logged_in():
		yield(highscore_api_service.register_highscore(_level_number, _moves, _timer), "finished")
	
	funcref(_callback_obj, _callback_func).call_func(true)
	_on_close()
