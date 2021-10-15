class_name LevelSelector
extends Control

const _items_in_row := 10


func _ready() -> void:
	var current_row_size := 0
	var current_hbox: HBoxContainer = null

	$CenterContainer.visible = false
	$BottomBox.visible = false

	var green_button_style = load("res://resources/green_button_style.tres")
	var button_style = load("res://resources/button_style.tres")
	var button_font = load("res://resources/button_font.tres")

	var my_scores = {}
	var nick_update_result: AsyncResult = null
	
	if auth_context.is_logged_in():
		nick_update_result = login_manager.refetch_user_data()
		var result: AsyncResult = yield(highscore_api_service.get_user_highscores(), "finished")
		if not result.is_error():
			my_scores = result.get_result()
	else:
		$CenterContainer/VBoxContainer/OfflineLabel.visible = true

	for i in range(constants.LEVEL_COUNT):
		if current_hbox == null:
			current_hbox = HBoxContainer.new()
			current_hbox.add_constant_override("separation", 10)

		var new_button := Button.new()
		new_button.text = str(i + 1)

		var current_style: StyleBox

		if my_scores.has(str(i + 1)):
			current_style = green_button_style
		else:
			current_style = button_style

		new_button.add_stylebox_override("hover", current_style)
		new_button.add_stylebox_override("pressed", current_style)
		new_button.add_stylebox_override("focus", current_style)
		new_button.add_stylebox_override("normal", current_style)
		new_button.add_font_override("font", button_font)
		new_button.add_color_override("font_color", "#556270")
		new_button.add_color_override("font_color_hover", "#4c5864")
		new_button.add_color_override("font_color_pressed", "#000000")
		new_button.mouse_default_cursor_shape = Control.CURSOR_POINTING_HAND
		new_button.rect_min_size = Vector2(47, 47)
		new_button.connect("pressed", self, "_on_button_pressed", [i + 1])
		current_hbox.add_child(new_button)

		current_row_size += 1

		if current_row_size == _items_in_row:
			current_row_size = 0
			$CenterContainer/VBoxContainer.add_child(current_hbox)
			current_hbox = null

	if current_hbox != null:
		$CenterContainer/VBoxContainer.add_child(current_hbox)

	if nick_update_result != null and not nick_update_result.is_finished():
		yield(nick_update_result, "finished")

	if auth_context.get_nickname() != "":
		$CenterContainer/VBoxContainer/Label.text += ", " + auth_context.get_nickname()

	$CenterContainer.visible = true
	$BottomBox.visible = true

	anim_helper.slide_in_from_down($BottomBox, 0.2)
	anim_helper.fade_in($CenterContainer/VBoxContainer, 0.2)


func animate_out() -> void:
	anim_helper.slide_out_down($BottomBox, 0.2)
	anim_helper.fade_out($CenterContainer/VBoxContainer, 0.2)


func _on_button_pressed(level: int) -> void:
	scene_switcher.change_scene("res://src/levels/level_" + str(level) + ".tscn", 0.2)


func _on_BackButton_pressed() -> void:
	scene_switcher.change_scene("res://src/ui/main_menu.tscn", 0.2)
