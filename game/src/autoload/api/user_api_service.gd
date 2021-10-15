class_name UserAPIService
extends APIServiceBase


func register_user(password: String, show_indicator: bool = true) -> AsyncResult:
	return _handle_json_results(
		http_handler.do_post(
			constants.get_base_url() + "user",
			JSON.print({"username": null, "password": password}),
			show_indicator
		)
	)


func change_nickname(nickname: String, show_indicator: bool = true) -> AsyncResult:
	return http_handler.do_patch_with_auth(
		constants.get_base_url() + "user/nickname", nickname, show_indicator
	)


func get_user_details(show_indicator: bool = true) -> AsyncResult:
	return _handle_json_results(
		http_handler.do_get_with_auth(constants.get_base_url() + "user", show_indicator)
	)


func login(username: String, password: String, show_indicator: bool = true) -> AsyncResult:
	return _handle_json_results(
		http_handler.do_post(
			constants.get_base_url() + "session",
			JSON.print({"username": username, "password": password}),
			show_indicator
		)
	)


func logout(show_indicator: bool = true) -> AsyncResult:
	return _handle_json_results(
		http_handler.do_delete_with_auth(constants.get_base_url() + "session", show_indicator)
	)
