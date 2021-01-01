.PHONY: swagger-codegen
swagger-codegen:
ifeq (,$(wildcard ./bin/swagger-codegen-cli.jar))
	@{ \
	set -e ;\
	mkdir -p bin ;\
	curl -sSLo ./bin/swagger-codegen-cli.jar https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/2.4.9/swagger-codegen-cli-2.4.9.jar ;\
	}
endif

.PHONY: model
model: swagger-codegen
	@echo Generating models...
	@echo
	@java -Dmodels -jar ./bin/swagger-codegen-cli.jar generate -i swagger/swagger.yaml -l go -c swagger/swagger.conf -t swagger/template/go -o ./pkg/model >/dev/null 2>&1

.PHONY: build
build: model
	@echo Building...
	@echo
	@go build -o bin/zundoko-client ./cmd

.PHONY: start-server
start-server:
ifeq (,$(wildcard ./server.pid))
	@{ \
	set -e ;\
	cd react-redux-scaffold ;\
	npm install ;\
	npm run mock & echo "$$!" > ../server.pid ;\
	}
endif

.PHONY: stop-server
stop-server:
ifneq (,$(wildcard ./server.pid))
	@{ \
	set -e ;\
	kill -INT $$(ps o pid,ppid,cmd | grep $$(cat server.pid) | grep node | awk '{print $$1}') ;\
	rm -f server.pid ;\
	}
endif

.PHONY: ginkgo
ginkgo:
ifeq (,$(shell which ginkgo 2>/dev/null))
	@{ \
	set -e ;\
	go get github.com/onsi/ginkgo/ginkgo ;\
	}
endif

.PHONY: test-suite
test-suite: ginkgo
ifeq (,$(pkg))
	@echo specify a package. e.g. pkg=pkg/runner
	@false
endif
ifeq (,$(wildcard $(pkg)))
	@echo package ${pkg} does not exist.
else
	@{ \
	set -e ;\
	cd ${pkg} ;\
	ginkgo bootstrap -internal ;\
	}
endif

.PHONY: test-template
test-template: ginkgo
ifeq (,$(src))
	@echo specify a source file. e.g. src=pkg/runner/runner.go
	@false
endif
ifeq (,$(wildcard $(src)))
	@echo source ${src} does not exist.
else
	@{ \
	set -e ;\
	cd $$(dirname ${src}) ;\
	ginkgo generate -internal $$(basename ${src}) ;\
	}
endif

.PHONY: mockgen
mockgen:
ifeq (,$(shell which mockgen 2>/dev/null))
	@{ \
	set -e ;\
	go get github.com/golang/mock/mockgen ;\
	}
endif

.PHONY: mock
mock: model mockgen
	@echo Generating mocks...
	@echo
	@for GO_FILE in $$(find ./pkg -name "*.go" -not -name "*_test.go" -not -name "doc.go"); do\
		MOCK_DIR=mock/$$(dirname $$(dirname $$GO_FILE))/mock_$$(basename $$(dirname $$GO_FILE)) ;\
		mkdir -p $${MOCK_DIR} ;\
		mockgen -source=$$GO_FILE -destination $${MOCK_DIR}/$$(basename $$GO_FILE) ;\
	done
	@# Remove emply mock files.
	@rm -f $$(find mock/ -name "*.go" | xargs grep -iL "func ")

.PHONY: test mock
test: ginkgo mock
	@echo Running unit tests...
	@echo
	@ginkgo -r -cover

.PHONY: clean
clean: stop-server
	@rm -rf pkg/model bin server.pid mock
