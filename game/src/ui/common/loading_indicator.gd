class_name LoadingIndicator
extends ColorRect


func _ready() -> void:
	loading_manager.connect("loading_changed", self, "_on_loading_changed")
	_on_loading_changed(loading_manager.is_loading())


func _on_loading_changed(loading: bool) -> void:
	if loading:
		if $Timer.is_stopped():
			$Timer.start()
	else:
		if not $Timer.is_stopped():
			$Timer.stop()

		self.mouse_filter = Control.MOUSE_FILTER_IGNORE
		self.visible = false


func _on_Timer_timeout() -> void:
	self.mouse_filter = Control.MOUSE_FILTER_STOP
	self.visible = true
