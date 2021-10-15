extends Node2D

const _left_dir := Vector2(-1, 0)
const _right_dir := Vector2(1, 0)
const _up_dir := Vector2(0, -1)
const _down_dir := Vector2(0, 1)

var _grid: Grid = null
var _moving: bool = false
var _movement_queue = []
var _move_towards_mouse: bool = false

signal player_moved


func _ready() -> void:
	_grid = $"../Grid"


func _get_my_pos() -> Vector2:
	return _grid.world_to_map(position)


func _move_to(direction: Vector2) -> void:
	_moving = true

	if direction == _left_dir:
		$AnimationPlayer.play("move_left")
	elif direction == _right_dir:
		$AnimationPlayer.play("move_right")
	elif direction == _down_dir:
		$AnimationPlayer.play("move_down")
	elif direction == _up_dir:
		$AnimationPlayer.play("move_up")

	$AnimationPlayer.advance(0.0)
	translate(direction * _grid.cell_size)
	emit_signal("player_moved")

	yield($AnimationPlayer, "animation_finished")
	yield(get_tree().create_timer(.1), "timeout")

	_moving = false


func _handle_button_press(action: String, direction: Vector2) -> void:
	if Input.is_action_pressed(action):
		if not _movement_queue.has(direction):
			_movement_queue.push_front(direction)
	else:
		if _movement_queue.has(direction):
			_movement_queue.erase(direction)


func _unhandled_input(event):
	if event is InputEventMouseButton and event.button_index == BUTTON_LEFT:
		_move_towards_mouse = event.pressed

	if event is InputEventScreenTouch:
		_move_towards_mouse = event.pressed


func _handle_mouse() -> void:
	var mouse_pos := _grid.world_to_map(_grid.get_local_mouse_position())
	var my_pos := _get_my_pos()

	var directions = []

	if mouse_pos.x > my_pos.x:
		directions.append(_right_dir)
	elif mouse_pos.x < my_pos.x:
		directions.append(_left_dir)

	if mouse_pos.y > my_pos.y:
		directions.append(_down_dir)
	elif mouse_pos.y < my_pos.y:
		directions.append(_up_dir)

	if abs(my_pos.y - mouse_pos.y) > abs(my_pos.x - mouse_pos.x):
		directions.invert()

	for direction in directions:
		if _grid.can_move_to(_get_my_pos(), _get_my_pos() + direction):
			_move_to(direction)
			return


func _process(delta: float) -> void:
	_handle_button_press("ui_up", _up_dir)
	_handle_button_press("ui_down", _down_dir)
	_handle_button_press("ui_left", _left_dir)
	_handle_button_press("ui_right", _right_dir)

	if _moving:
		return

	if _move_towards_mouse:
		_handle_mouse()

	for direction in _movement_queue:
		if _grid.can_move_to(_get_my_pos(), _get_my_pos() + direction):
			_move_to(direction)
			return
