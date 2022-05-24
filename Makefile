
export GO111MODULE=on

#name openebsctl to be a kubectl-plugin
OPENEBSCTL=kubectl-openebs

.PHONY: test
test:
	go test ./pkg/... ./cmd/... -coverprofile cover.out

.PHONY: bin
bin: fmt vet
	go build -o bin/graviton-scheduler-extender github.com/marccampbell/graviton-scheduler-extender/cmd/graviton-scheduler-extender

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: openebsctl
openebsctl:
	@echo "----------------------------"
	@echo "--> openebsctl                    "
	@echo "----------------------------"
	@PNAME=OPENEBSCTL CTLNAME=${OPENEBSCTL} sh -c "'$(PWD)/scripts/build.sh'"
	@echo "--> Removing old directory..."
	@sudo rm -rf /usr/local/bin/${OPENEBSCTL}
	@echo "----------------------------"
	@echo "copying new openebsctl"
	@echo "----------------------------"
	@sudo mkdir -p  /usr/local/bin/
	@sudo cp -a "$(PWD)/bin/OPENEBSCTL/darwin_arm64/${OPENEBSCTL}"  /usr/local/bin/${OPENEBSCTL}
	@echo "=> copied to /usr/local/bin"

.PHONY: run
run:
	kubectl delete -f ./install/graviton-scheduler-extender.yaml || true
	docker build -t ttl.sh/graviton-scheduler-extender:24h .
	docker push ttl.sh/graviton-scheduler-extender:24h
	kubectl apply -f ./install/graviton-scheduler-extender.yaml