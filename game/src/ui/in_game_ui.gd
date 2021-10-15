class_name InGameUI
extends Control

export(int) var level_number = 0

var _grid: Grid = null
var _game: Node2D = null
var _start_timestamp
var _timer: float = 0.0
var _moves: int = 0
var _stop_timer: bool = false

var _previous_best_moves := 0
var _previous_best_time := 0


func _ready() -> void:
	_grid = $"../Game/Grid"
	_game = $"../Game"
	get_tree().root.connect("size_changed", self, "_on_viewport_size_changed")
	_on_viewport_size_changed()
	var label := $MarginContainer/HBoxContainer/VBoxContainer2/LevelLabel
	label.text = "Level: " + str(level_number)
	$"../Game/Player".connect("player_moved", self, "_on_player_moved")

	yield(get_tree().create_timer(0), "timeout")
	_start_timestamp = OS.get_system_time_secs()

	anim_helper.slide_in_from_up($MarginContainer, 0.2)
	anim_helper.slide_in_from_down($MarginContainer2, 0.2)

	if auth_context.is_logged_in():
		var result: AsyncResult = yield(highscore_api_service.get_user_highscores(), "finished")

		if not result.is_error():
			var scores = result.get_result()
			if scores.has(str(level_number)):
				_previous_best_moves = scores[str(level_number)].moves
				_previous_best_time = scores[str(level_number)].completionTime


func animate_out() -> void:
	anim_helper.slide_out_up($MarginContainer, 0.2)
	anim_helper.slide_out_down($MarginContainer2, 0.2)


func _on_player_moved() -> void:
	_moves = _moves + 1
	$MarginContainer/HBoxContainer/VBoxContainer3/MovesLabel.text = "Moves: " + str(_moves)

	if _grid.is_completed():
		_stop_timer = true
		dialog_opener.open_custom_dialog(
			load("res://src/ui/level_complete_dialog.tscn"),
			self,
			"_on_complete_response",
			[
				level_number,
				_moves,
				_timer,
				_get_formatted_timer(_timer),
				_previous_best_moves,
				_previous_best_time,
				_get_formatted_timer(_previous_best_time)
			]
		)


func _on_complete_response(response: bool) -> void:
	if response and level_number < constants.LEVEL_COUNT:
		scene_switcher.change_scene(
			"res://src/levels/level_" + str(level_number + 1) + ".tscn", 0.2
		)
	else:
		scene_switcher.change_scene("res://src/ui/level_selector.tscn", 0.2)


func _get_formatted_timer(timer: float) -> String:
	var hours: float = floor(timer / (60 * 60))
	var minutes: float = floor((timer - hours * 60 * 60) / 60)
	var seconds: float = timer - hours * 60 * 60 - minutes * 60

	return "%02d:%02d:%02d" % [hours, minutes, seconds]


func _process(delta):
	if _stop_timer:
		return

	_timer = OS.get_system_time_secs() - _start_timestamp

	var text = _get_formatted_timer(_timer)

	$MarginContainer/HBoxContainer/VBoxContainer3/CenterContainer/TimerLabel.text = text


func _get_level_size() -> Vector2:
	var highest := Vector2(0, 0)

	for v in _grid.get_used_cells():
		if v.x > highest.x:
			highest.x = v.x

		if v.y > highest.y:
			highest.y = v.y

	return (highest + Vector2(1, 1)) * _grid.cell_size


func _on_viewport_size_changed() -> void:
	var max_size := get_viewport().get_visible_rect().size * .9
	var level_size := _get_level_size()

	var scale := max_size / level_size

	if scale.x > 1.0:
		scale.x = 1.0

	if scale.y > 1.0:
		scale.y = 1.0

	if scale.x < scale.y:
		scale.y = scale.x
	elif scale.y < scale.x:
		scale.x = scale.y

	_game.scale = scale
	_game.position = (
		get_viewport().get_visible_rect().size / 2
		- (_get_level_size() / 2 * scale)
	)


func _on_BackButton_pressed() -> void:
	dialog_opener.open_question_dialog(
		"Are you sure",
		"Are you sure you want to go back to the level selector? Your progress will be lost.",
		self,
		"_on_back_response"
	)


func _on_back_response(response: bool) -> void:
	if response:
		scene_switcher.change_scene("res://src/ui/level_selector.tscn", 0.2)


func _on_RestartButton_pressed():
	dialog_opener.open_question_dialog(
		"Are you sure",
		"Are you sure you want to reset the level? Your progress will be lost.",
		self,
		"_on_reset_response"
	)


func _on_reset_response(response: bool) -> void:
	if response:
		scene_switcher.change_scene(
			"res://src/levels/level_" + str(level_number) + ".tscn", 0.2
		)
