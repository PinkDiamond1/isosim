HTTP_PORT=8088
LOG_LEVEL=INFO

build: build-isosim build-isoserver
build-isosim:
	@go build -o ./bin/isosim ./cmd/isosim
	@echo "✅ isosim done"

build-isoserver:
	@go build -o ./bin/isoserver ./cmd/isoserver
	@echo "✅ isoserver done"

run-isosim:
	@echo "🌏 http://localhost:${HTTP_PORT}/"
	@./bin/isosim \
	    -http-port ${HTTP_PORT} --log-level ${LOG_LEVEL} \
	    -specs-dir ${PWD}/test/testdata/specs \
		-html-dir ${PWD}/web \
		-data-dir ${PWD}/test/testdata/appdata

run-isoserver:
	@echo "TODO"
	@./bin/isoserver \
	    -specs-dir ${PWD}/test/testdata/specs \
		--def-file ${PWD}/test/testdata/appdata/2/IsoMiniSpec_Server_01.srvdef.json


clean:
	@rm -rf ./bin
