class_name AnimHelper
extends Node

var _tween: Tween = Tween.new()


func _ready():
	add_child(_tween)


func slide_out_up(node: Control, time = 0.2) -> void:
	_tween.interpolate_property(
		node,
		"rect_position",
		node.rect_position,
		Vector2(node.rect_position.x, -node.rect_size.y),
		time,
		Tween.TRANS_QUAD,
		Tween.EASE_OUT
	)
	_tween.start()


func slide_out_down(node: Control, time = 0.2) -> void:
	_tween.interpolate_property(
		node,
		"rect_position",
		node.rect_position,
		Vector2(node.rect_position.x, node.get_parent().rect_size.y + 1),
		time,
		Tween.TRANS_QUAD,
		Tween.EASE_OUT
	)
	_tween.start()


func slide_in_from_up(node: Control, time = 0.2) -> void:
	var orig_pos: Vector2 = node.rect_position
	node.rect_position = Vector2(orig_pos.x, -node.rect_size.y)

	_tween.interpolate_property(
		node,
		"rect_position",
		node.rect_position,
		orig_pos,
		time,
		Tween.TRANS_QUAD,
		Tween.EASE_OUT
	)
	_tween.start()


func slide_canvasitem_in_from_up(node: CanvasItem, height: int, time: float = 0.2) -> void:
	var orig_pos: Vector2 = node.position
	node.position = Vector2(orig_pos.x, orig_pos.y - height)

	_tween.interpolate_property(
		node, "position", node.position, orig_pos, time, Tween.TRANS_QUAD, Tween.EASE_OUT
	)
	_tween.start()


func slide_in_from_down(node: Control, time = 0.2) -> void:
	var orig_pos: Vector2 = node.rect_position
	node.rect_position = Vector2(orig_pos.x, node.get_parent().rect_size.y + 1)

	_tween.interpolate_property(
		node,
		"rect_position",
		node.rect_position,
		orig_pos,
		0.2,
		Tween.TRANS_QUAD,
		Tween.EASE_OUT
	)
	_tween.start()


func fade_out(node: CanvasItem, time = 0.2) -> void:
	_tween.interpolate_property(
		node,
		"modulate",
		node.modulate,
		Color(1, 1, 1, 0),
		time,
		Tween.TRANS_QUAD,
		Tween.EASE_OUT
	)
	_tween.start()


func fade_in(node: CanvasItem, time = 0.2) -> void:
	node.modulate = Color(1, 1, 1, 0)
	_tween.interpolate_property(
		node,
		"modulate",
		node.modulate,
		Color(1, 1, 1, 1),
		time,
		Tween.TRANS_QUAD,
		Tween.EASE_OUT
	)
	_tween.start()


func scale_in(node: Node2D, time = 0.2) -> void:
	var orig_scale: Vector2 = node.scale
	node.scale = Vector2(0, 0)
	_tween.interpolate_property(
		node, "scale", Vector2(0, 0), orig_scale, time, Tween.TRANS_QUAD, Tween.EASE_OUT
	)
	_tween.start()


func scale_to(node: Node2D, target: Vector2, time = 0.2) -> void:
	_tween.interpolate_property(
		node, "scale", node.scale, target, time, Tween.TRANS_QUAD, Tween.EASE_OUT
	)
	_tween.start()


func zoom_to(camera: Camera2D, target: Vector2, time: float = 0.2) -> void:
	_tween.interpolate_property(
		camera, "zoom", camera.zoom, target, time, Tween.TRANS_QUAD, Tween.EASE_OUT
	)
	_tween.start()
