class_name DialogOpener
extends Node


func _ready() -> void:
	pass


func _create_dialog(title: String, description: String) -> GeneralDialog:
	var holder: CanvasLayer = CanvasLayer.new()
	var dialog: GeneralDialog = load("res://src/ui/general_dialog.tscn").instance()
	dialog.set_holder(holder)
	dialog.set_title(title)
	dialog.set_description(description)
	get_tree().get_root().add_child(holder)
	holder.add_child(dialog)
	dialog.connect("close", self, "_on_dialog_closed")
	return dialog


func show_toast_message(text: String) -> ToastDialog:
	var holder: CanvasLayer = CanvasLayer.new()
	var dialog: ToastDialog = load("res://src/ui/toast_dialog.tscn").instance()
	dialog.set_holder(holder)
	dialog.set_description(text)
	get_tree().get_root().add_child(holder)
	holder.add_child(dialog)
	dialog.connect("close", self, "_on_dialog_closed")
	return dialog


func open_custom_dialog(
	dialog_scene, callback_obj = null, callback_func: String = "", arguments = []
) -> void:
	var holder: CanvasLayer = CanvasLayer.new()
	var dialog = dialog_scene.instance()
	dialog.set_holder(holder)

	if callback_obj != null:
		dialog.set_callback(callback_obj, callback_func)

	if arguments.size() > 0:
		dialog.set_arguments(arguments)

	get_tree().get_root().add_child(holder)
	holder.add_child(dialog)
	dialog.connect("close", self, "_on_custom_dialog_closed")


func open_question_dialog(
	title: String, description: String, callback_obj, callback_func: String
) -> void:
	var dialog: GeneralDialog = _create_dialog(title, description)
	dialog.set_question_mode(callback_obj, callback_func)


func open_error_dialog(title: String, description: String) -> void:
	var dialog: GeneralDialog = _create_dialog(title, description)
	dialog.set_warning_mode()


func _on_custom_dialog_closed(dialog, holder) -> void:
	holder.remove_child(dialog)


func _on_dialog_closed(dialog: GeneralDialog, holder) -> void:
	holder.remove_child(dialog)
	get_tree().get_root().remove_child(holder)
