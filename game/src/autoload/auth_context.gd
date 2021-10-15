class_name AuthContext
extends Node

var _token: String = ""
var _username: String = ""
var _nickname: String = ""


func is_logged_in() -> bool:
	return _token != ""


func get_nickname() -> String:
	return _nickname


func get_username() -> String:
	return _username


func update_user_data(username: String, nickname: String) -> void:
	_nickname = nickname
	_username = username


func update_token(token: String) -> void:
	_token = token


func clear_token() -> void:
	_token = ""


func get_token_header() -> PoolStringArray:
	var ret: PoolStringArray
	ret.append("X-Auth-Token: %s" % _token)
	return ret
