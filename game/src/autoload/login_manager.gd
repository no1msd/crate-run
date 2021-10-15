class_name LoginManager
extends Node


func _get_save_file_name():
	return "user://login.save"


func _has_saved_login():
	var data = File.new()
	return data.file_exists(_get_save_file_name())


func _get_saved_login():
	var data = File.new()
	data.open(_get_save_file_name(), File.READ)
	var login = data.get_var()
	data.close()
	return login


func _save_login(login):
	var data = File.new()
	data.open(_get_save_file_name(), File.WRITE)
	data.store_var(login)
	data.close()


func _generate_random_password() -> String:
	var chars = [
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"q",
		"w",
		"e",
		"r",
		"t",
		"z",
		"u",
		"i",
		"o",
		"p",
		"a",
		"s",
		"d",
		"f",
		"g",
		"h",
		"j",
		"k",
		"l",
		"y",
		"x",
		"c",
		"v",
		"b",
		"n",
		"m",
		"Q",
		"W",
		"E",
		"R",
		"T",
		"Z",
		"U",
		"I",
		"O",
		"P",
		"A",
		"S",
		"D",
		"F",
		"G",
		"H",
		"J",
		"K",
		"L",
		"Y",
		"X",
		"C",
		"V",
		"B",
		"N",
		"M",
		",",
		".",
		";",
		"!",
		"%",
		"=",
		"#",
		"@",
		"&",
		"[",
		"]"
	]

	var rng: RandomNumberGenerator = RandomNumberGenerator.new()
	rng.randomize()

	var password: String = ""

	for i in range(16):
		password += chars[rng.randi_range(0, chars.size() - 1)]

	return password


func _on_refetch_task_finished(api_result: AsyncResult, login_result: AsyncResult) -> void:
	if api_result.is_error():
		login_result.failure(api_result.get_error())
	else:
		var user = api_result.get_result()
		auth_context.update_user_data(user.username, user.nickname)
		login_result.success({})


func refetch_user_data() -> AsyncResult:
	var login_result := AsyncResult.new()
	var api_result := user_api_service.get_user_details()
	api_result.connect("finished", self, "_on_refetch_task_finished", [login_result])
	return login_result


func _login(login_result: AsyncResult) -> void:
	var data = JSON.parse(_get_saved_login()).result

	var result: AsyncResult = yield(
		user_api_service.login(data.username, data.password), "finished"
	)

	if not result.is_error():
		var session = result.get_result()
		auth_context.update_token(session.token)
		auth_context.update_user_data(session.user.username, session.user.nickname)

		login_result.success({})
	else:
		auth_context.clear_token()
		_register(login_result)


func _register(login_result: AsyncResult) -> void:
	var generated_password = _generate_random_password()

	var result: AsyncResult = yield(
		user_api_service.register_user(generated_password), "finished"
	)

	if not result.is_error():
		var session = result.get_result()
		auth_context.update_token(session.token)
		auth_context.update_user_data(session.user.username, session.user.nickname)
		_save_login(
			JSON.print({"username": session.user.username, "password": generated_password})
		)

		login_result.success({})
	else:
		auth_context.clear_token()
		login_result.failure(result.get_error())


func login_or_register() -> AsyncResult:
	var result := AsyncResult.new()

	if _has_saved_login():
		_login(result)
	else:
		_register(result)

	return result
