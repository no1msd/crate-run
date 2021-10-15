class_name Grid
extends TileMap

var _crates = []


func _ready() -> void:
	for crate in $"../Crates".get_children():
		_crates.append(crate)
		crate.update_is_target(_is_target(world_to_map(crate.position)))


func _is_wall(pos: Vector2) -> bool:
	var cellv = get_cellv(pos)

	if cellv == INVALID_CELL:
		return false

	return tile_set.tile_get_name(cellv) == "wall"


func _is_target(pos: Vector2) -> bool:
	var cellv = get_cellv(pos)

	if cellv == INVALID_CELL:
		return false

	return tile_set.tile_get_name(cellv) == "target"


func is_completed() -> bool:
	for crate in _crates:
		if not _is_target(world_to_map(crate.position)):
			return false

	return true


func can_move_to(from: Vector2, to: Vector2) -> bool:
	if _is_wall(to):
		return false

	var pushed_crate: Crate = null
	var crate_or_wall_after_next: bool = false
	var next_pos = to + (to - from)

	for crate in _crates:
		var crate_pos = world_to_map(crate.position)

		if crate_pos == to:
			pushed_crate = crate

		if crate_pos == next_pos or _is_wall(next_pos):
			crate_or_wall_after_next = true

	if pushed_crate and not crate_or_wall_after_next:
		pushed_crate.move(cell_size, to - from)
		pushed_crate.update_is_target(_is_target(next_pos))

	if pushed_crate and crate_or_wall_after_next:
		pushed_crate.play_stuck_animation(to - from)
		return false

	return true
