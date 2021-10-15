class_name HighScoreAPIService
extends APIServiceBase


func get_global_highscores(show_indicator: bool = true) -> AsyncResult:
	return _handle_json_results(
		http_handler.do_get(constants.get_base_url() + "highscores", show_indicator)
	)


func get_level_highscores(number: int, show_indicator: bool = true) -> AsyncResult:
	return _handle_json_results(
		http_handler.do_get(
			constants.get_base_url() + "highscores/level/" + String(number), show_indicator
		)
	)


func get_user_highscores(show_indicator: bool = true) -> AsyncResult:
	return _handle_json_results(
		http_handler.do_get_with_auth(
			constants.get_base_url() + "highscores/user", show_indicator
		)
	)


func register_highscore(
	number: int, moves: int, completion_time: int, show_indicator: bool = true
) -> AsyncResult:
	return http_handler.do_put_with_auth(
		constants.get_base_url() + "highscores/level/" + String(number),
		JSON.print({"moves": moves, "completionTime": completion_time}),
		show_indicator
	)
