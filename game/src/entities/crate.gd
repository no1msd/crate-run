class_name Crate
extends Node2D

var _was_on_target: bool = false


func play_stuck_animation(direction: Vector2) -> void:
	if direction == Vector2(1, 0):
		$StuckAnimationPlayer.play("Stuck Right")
	if direction == Vector2(-1, 0):
		$StuckAnimationPlayer.play("Stuck Left")
	if direction == Vector2(0, 1):
		$StuckAnimationPlayer.play("Stuck Down")
	if direction == Vector2(0, -1):
		$StuckAnimationPlayer.play("Stuck Up")


func update_is_target(is_target: bool) -> void:
	if _was_on_target and not is_target:
		$TickAnimationPlayer.play("fade_out")
		$TickAnimationPlayer.advance(0.0)
		yield($TickAnimationPlayer, "animation_finished")
		_was_on_target = false
		$IdleAnimationPlayer.play("idle")

	if not _was_on_target and is_target:
		_was_on_target = true
		$TickAnimationPlayer.play("fade_in")
		$TickAnimationPlayer.advance(0.0)
		yield(get_tree().create_timer(0), "timeout")
		$IdleAnimationPlayer.stop()
		$Tween.interpolate_property(
			$IdleHolder,
			"scale",
			$IdleHolder.scale,
			Vector2(0.9, 0.9),
			0.1,
			Tween.TRANS_QUAD,
			Tween.EASE_OUT
		)
		$Tween.start()


func move(cell_size: Vector2, direction: Vector2) -> void:
	if direction == Vector2(-1, 0):
		$AnimationPlayer.play("move_left")
	elif direction == Vector2(1, 0):
		$AnimationPlayer.play("move_right")
	elif direction == Vector2(0, 1):
		$AnimationPlayer.play("move_down")
	elif direction == Vector2(0, -1):
		$AnimationPlayer.play("move_up")

	$AnimationPlayer.advance(0.0)
	translate(direction * cell_size)
