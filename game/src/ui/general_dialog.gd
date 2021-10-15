class_name GeneralDialog
extends ColorRect

signal close(dialog, holder)

var _holder: CanvasLayer = null

var _callback_obj = null
var _callback_func: String = ""


func _ready():
	anim_helper.fade_in(self)
	anim_helper.slide_in_from_down($Control)
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/OKButton.connect(
		"button_up", self, "_on_ok"
	)
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/NoButton.connect(
		"button_up", self, "_on_no"
	)
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/YesButton.connect(
		"button_up", self, "_on_yes"
	)


func set_holder(holder: CanvasLayer) -> void:
	_holder = holder


func set_title(title: String) -> void:
	pass
	#$Control/MarginContainer/MarginContainer/VBoxContainer/Label.text = title


func set_description(description: String) -> void:
	$Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/HBoxContainer/Label3.text = description


func set_warning_mode() -> void:
	$Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/HBoxContainer/CenterContainer/InfoIcon.visible = false
	$Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/HBoxContainer/WarningIcon.visible = true
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/OKButton.visible = true
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/NoButton.visible = false
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/YesButton.visible = false


func set_info_mode() -> void:
	$Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/HBoxContainer/CenterContainer/InfoIcon.visible = true
	$Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/HBoxContainer/WarningIcon.visible = false
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/OKButton.visible = true
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/NoButton.visible = false
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/YesButton.visible = false


func set_question_mode(callback_obj, callback_func: String) -> void:
	$Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/HBoxContainer/CenterContainer/InfoIcon.visible = true
	$Control/MarginContainer/MarginContainer/VBoxContainer/MarginContainer/HBoxContainer/WarningIcon.visible = false
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/OKButton.visible = false
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/NoButton.visible = true
	$Control/MarginContainer/MarginContainer/VBoxContainer/HBoxContainer/YesButton.visible = true
	_callback_obj = callback_obj
	_callback_func = callback_func


func _on_ok() -> void:
	_on_close()


func _on_yes() -> void:
	if _callback_obj != null:
		funcref(_callback_obj, _callback_func).call_func(true)
	_on_close()


func _on_no() -> void:
	if _callback_obj != null:
		funcref(_callback_obj, _callback_func).call_func(false)
	_on_close()


func _on_close() -> void:
	anim_helper.fade_out(self)
	anim_helper.slide_out_down($Control)

	yield(get_tree().create_timer(0.2), "timeout")

	emit_signal("close", self, _holder)
