HTTP_PORT=8088
LOG_LEVEL=INFO

build: build-isosim
build-isosim:
	@go build -o ./bin/isosim ./cmd/isosim
	@echo "✅ isosim done"

build-isoserver:
	@go build -o ./bin/isoserver ./cmd/isoserver
	@echo "✅ isoserver done"

run-isosim:
	@echo "🌏 http://localhost:${HTTP_PORT}/"
	@./bin/isosim -http-port ${HTTP_PORT} --log-level ${LOG_LEVEL} \
	              -specs-dir ${PWD}/test/testdata/specs \
				  -html-dir ${PWD}/web \
				  -data-dir ${PWD}/test/testdata/appdata

run-isoserver:
	@echo "TODO"

clean:
	@rm -rf ./bin
