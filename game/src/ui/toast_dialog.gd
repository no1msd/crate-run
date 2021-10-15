class_name ToastDialog
extends ColorRect

signal close(dialog, holder)

var _holder: CanvasLayer = null


func _ready():
	anim_helper.fade_in(self)

	yield(get_tree().create_timer(3.0), "timeout")

	_on_close()


func set_holder(holder: CanvasLayer) -> void:
	_holder = holder


func set_description(description: String) -> void:
	$Control/MarginContainer/MarginContainer/HBoxContainer/Label3.text = description


func _on_close() -> void:
	anim_helper.fade_out(self)

	yield(get_tree().create_timer(0.2), "timeout")

	emit_signal("close", self, _holder)
