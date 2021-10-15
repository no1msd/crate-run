class_name APIServiceBase
extends Node


func _on_json_http_finished(http_result: AsyncResult, api_result: AsyncResult) -> void:
	var response = http_result.get_result()

	if response.result == 0 and response.response_code == 200:
		var body = response.body.get_string_from_utf8()
		if body != "":
			api_result.success(JSON.parse(body).result)
		else:
			api_result.success({})
	else:
		api_result.failure(response.response_code)


func _on_binary_http_finished(http_result: AsyncResult, api_result: AsyncResult) -> void:
	var response = http_result.get_result()

	if response.result == 0 and response.response_code == 200:
		api_result.success(response.body)
	else:
		api_result.failure(response.response_code)


func _handle_json_results(http_result: AsyncResult) -> AsyncResult:
	var api_result := AsyncResult.new()
	http_result.connect("finished", self, "_on_json_http_finished", [api_result])
	return api_result


func _handle_binary_results(http_result: AsyncResult) -> AsyncResult:
	var api_result := AsyncResult.new()
	http_result.connect("finished", self, "_on_binary_http_finished", [api_result])
	return api_result
