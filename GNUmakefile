default: testacc

# Run acceptance tests
.PHONY: testacc

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

download-schema:
	curl https://raw.githubusercontent.com/linear/linear/master/packages/sdk/src/schema.graphql > schema.graphql
