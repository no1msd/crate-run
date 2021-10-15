class_name HTTPHandler
extends Node

const MAX_RUNNING_REQUESTS: int = 8

var _current_running_requests: int = 0

var _queued_requests = []


func do_patch(url: String, data: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(url, [], HTTPClient.METHOD_PATCH, data, null, show_indicator)


func do_patch_with_auth(url: String, data: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(
		url,
		auth_context.get_token_header(),
		HTTPClient.METHOD_PATCH,
		data,
		null,
		show_indicator
	)


func do_post(url: String, data: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(url, [], HTTPClient.METHOD_POST, data, null, show_indicator)


func do_post_with_auth(url: String, data: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(
		url, auth_context.get_token_header(), HTTPClient.METHOD_POST, data, null, show_indicator
	)


func do_put(url: String, data: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(url, [], HTTPClient.METHOD_PUT, data, null, show_indicator)


func do_put_with_auth(url: String, data: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(
		url, auth_context.get_token_header(), HTTPClient.METHOD_PUT, data, null, show_indicator
	)


func do_get(url: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(url, [], HTTPClient.METHOD_GET, "", null, show_indicator)


func do_get_with_auth(url: String, show_indicator: bool) -> AsyncResult:
	return _do_request_base(
		url, auth_context.get_token_header(), HTTPClient.METHOD_GET, "", null, show_indicator
	)


func _do_request_base(
	url: String,
	custom_headers: PoolStringArray,
	method,
	request_data: String,
	async_result: AsyncResult,
	show_indicator: bool
) -> AsyncResult:
	if async_result == null:
		async_result = AsyncResult.new()

	if _current_running_requests >= MAX_RUNNING_REQUESTS:
		_queued_requests.append(
			{
				"url": url,
				"custom_headers": custom_headers,
				"method": method,
				"request_data": request_data,
				"async_result": async_result,
				"show_indicator": show_indicator
			}
		)
		return async_result

	if show_indicator:
		loading_manager.start_loading()

	var http: HTTPRequest = HTTPRequest.new()
	add_child(http)

	if not OS.has_feature("JavaScript"):
		http.set_use_threads(true)

	http.request(url, custom_headers, false, method, request_data)
	_current_running_requests += 1

	http.connect(
		"request_completed", self, "_on_request_completed", [http, async_result, show_indicator]
	)

	return async_result


func _on_request_completed(
	result: int,
	response_code: int,
	headers: PoolStringArray,
	body: PoolByteArray,
	http: HTTPRequest,
	async_result: AsyncResult,
	show_indicator: bool
) -> void:
	if show_indicator:
		loading_manager.stop_loading()

	async_result.success(
		{"result": result, "response_code": response_code, "headers": headers, "body": body}
	)

	remove_child(http)

	_current_running_requests -= 1

	if _queued_requests.size() == 0:
		return

	var request = _queued_requests[0]
	_queued_requests.remove(0)

	call_deferred(
		"_do_request_base",
		request.url,
		request.custom_headers,
		request.method,
		request.request_data,
		request.async_result,
		request.show_indicator
	)
