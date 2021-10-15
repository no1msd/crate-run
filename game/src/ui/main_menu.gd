class_name MainMenu
extends Control


func _ready() -> void:
	anim_helper.slide_in_from_up($CenterContainer/VBoxContainer/TextureRect, 0.2)
	anim_helper.slide_in_from_down($CenterContainer/VBoxContainer/VBoxContainer2, 0.2)
	anim_helper.fade_in($CenterContainer/VBoxContainer/VBoxContainer/LineEdit, 0.2)
	anim_helper.fade_in($CenterContainer/VBoxContainer/TextureRect, 0.2)
	anim_helper.fade_in($CenterContainer/VBoxContainer/VBoxContainer2/Button, 0.2)
	anim_helper.fade_in(
		$CenterContainer/VBoxContainer/VBoxContainer2/CenterContainer/Button2, 0.2
	)

	var result: AsyncResult = yield(login_manager.login_or_register(), "finished")
	$CenterContainer/VBoxContainer/VBoxContainer2/Button.disabled = false
	if not result.is_error():
		$CenterContainer/VBoxContainer/VBoxContainer/LineEdit.text = auth_context.get_nickname()
	
	if not auth_context.is_logged_in():
		$CenterContainer/VBoxContainer/VBoxContainer/LineEdit.text = "Offline"
		$CenterContainer/VBoxContainer/VBoxContainer/LineEdit.editable = false
		$CenterContainer/VBoxContainer/VBoxContainer2/CenterContainer/Button2.disabled = true


func _on_Button_pressed() -> void:
	if auth_context.is_logged_in():
		var nickname: String = $CenterContainer/VBoxContainer/VBoxContainer/LineEdit.text
		var result: AsyncResult = yield(user_api_service.change_nickname(nickname), "finished")
	
	scene_switcher.change_scene("res://src/ui/level_selector.tscn", 0.2)


func animate_out() -> void:
	anim_helper.slide_out_up($CenterContainer/VBoxContainer/TextureRect, 0.2)
	anim_helper.slide_out_down($CenterContainer/VBoxContainer/VBoxContainer2, 0.2)
	anim_helper.fade_out($CenterContainer/VBoxContainer/VBoxContainer/LineEdit, 0.2)
	anim_helper.fade_out($CenterContainer/VBoxContainer/TextureRect, 0.2)
	anim_helper.fade_out($CenterContainer/VBoxContainer/VBoxContainer2/Button, 0.2)
	anim_helper.fade_out(
		$CenterContainer/VBoxContainer/VBoxContainer2/CenterContainer/Button2, 0.2
	)


func _on_LineEdit_focus_entered():
	if not OS.has_feature("JavaScript"):
		return

	var is_touch_js := "(('ontouchstart' in window) || (navigator.maxTouchPoints > 0) || (navigator.msMaxTouchPoints > 0))"

	if not JavaScript.eval(is_touch_js, true):
		return

	var result = JavaScript.eval(
		(
			"prompt('%s', '%s');"
			% ["Nickname:", $CenterContainer/VBoxContainer/VBoxContainer/LineEdit.text]
		),
		true
	)

	if result != null:
		$CenterContainer/VBoxContainer/VBoxContainer/LineEdit.text = result

	$CenterContainer/VBoxContainer/VBoxContainer/LineEdit.release_focus()


func _on_Button2_pressed() -> void:
	scene_switcher.change_scene("res://src/ui/highscores.tscn", 0.2)
