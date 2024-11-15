default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: start-pinot
start-pinot:
	@echo "Starting Pinot"
	@docker compose up -d pinot-controller pinot-server pinot-broker pinot-zookeeper

.PHONY: start-redpanda
start-redpanda:
	@echo "Starting Redpanda"
	@docker compose up -d redpanda-0 redpanda-console

.PHONY: create-dummy-topic
create-dummy-topic:
	@echo "Creating dummy topic"
	@rpk topic create dummy-topic -X brokers=localhost:19092

.PHONY: install
install:
	@echo "Installing Terraform"
	@go install

