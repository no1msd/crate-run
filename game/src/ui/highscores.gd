class_name HighScores
extends Control


func _add_row(place: int, name: String, score: int) -> void:
	var new_label = Label.new()
	new_label.text = str(place) + ". " + name + ": " + str(score)
	new_label.add_font_override("font", load("res://resources/fonts/ui_font_small.tres"))
	$CenterContainer/VBoxContainer/Top10Container.add_child(new_label)


func _ready() -> void:
	anim_helper.slide_in_from_down($BottomBox, 0.2)
	$CenterContainer.visible = false

	var result: AsyncResult = yield(highscore_api_service.get_global_highscores(), "finished")

	var count := 0
	var my_place := -1
	var my_score := 0

	if not result.is_error():
		var scores = result.get_result()
		for entry in scores:
			if entry.user.nickname == "":
				continue

			if entry.user.username == auth_context.get_username():
				my_place = count + 1
				my_score = entry.totalScore

			if count < 10:
				_add_row(count + 1, entry.user.nickname, entry.totalScore)
				count += 1

	if my_place >= 10:
		var text = str(my_place) + ". " + auth_context.get_nickname() + ": " + str(my_score)
		$CenterContainer/VBoxContainer/CurrentUserLabel.text = text
		$CenterContainer/VBoxContainer/CurrentUserMargin.visible = true
		$CenterContainer/VBoxContainer/CurrentUserLabel.visible = true

	$CenterContainer.visible = true
	anim_helper.fade_in($CenterContainer, 0.2)


func animate_out() -> void:
	anim_helper.slide_out_down($BottomBox, 0.2)
	anim_helper.fade_out($CenterContainer, 0.2)


func _on_BackButton_pressed() -> void:
	scene_switcher.change_scene("res://src/ui/main_menu.tscn", 0.2)
