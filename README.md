# crate.run - simple sokoban HTML5 game

![crate.run logo](https://raw.githubusercontent.com/no1msd/crate-run/main/game/resources/logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/no1msd/crate-run)](https://goreportcard.com/report/github.com/no1msd/crate-run)
[![Open Source Love](https://badges.frapsoft.com/os/mit/mit.svg?v=102)](https://github.com/ellerbrock/open-source-badge/)
[![Godot v3.4](https://img.shields.io/badge/Godot-v3.4-%23478cbf?logo=godot-engine&logoColor=white)](https://godotengine.org/)
[![Go v1.17](https://img.shields.io/badge/Go-v1.17-%23478cbf?logo=go&logoColor=white)](https://golang.org/)

crate.run is a simple HTML5 sokoban game built with the Godot game engine and Golang. Try it out 
online at [crate.run](https://crate.run)!

## Run with Docker Compose

```shell script
docker-compose up
```

This will export the Godot project to HTML5, build the API and run the resulting Docker image 
that serves everything through port 9080. After it's finished building you can play the game at 
http://localhost:9080.

If you want to change the URL used to access the game, change the `BASE_URL` environment 
variable in the `docker-compose.yml` file.

## Run API tests

```shell script
cd api
go test -v ./...
```

## License

crate.run is licensed under the [MIT license](https://github.com/no1msd/crate-run/blob/main/LICENSE).
