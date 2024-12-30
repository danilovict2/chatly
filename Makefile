run: build
	@templ generate
	@./bin/app
build:
	@go build -o bin/app main.go
css: 
	@npx tailwindcss -i ./views/css/app.css -o ./public/index.css --watch