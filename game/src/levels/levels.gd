class_name Level
extends Control


func _ready():
	anim_helper.fade_in($Game, 0.2)


func animate_out():
	$InGameUI.animate_out()
	anim_helper.fade_out($Game, 0.2)
