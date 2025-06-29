.PHONY: test coverage coverage-html

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

coverage-check:
	@go test -coverprofile=coverage.out ./... > /dev/null
	@coverage=$$(go tool cover -func=coverage.out | grep "total:" | grep -oE '[0-9]+\.[0-9]+'); \
	if [ $$(echo "$$coverage < 80" | bc) -eq 1 ]; then \
		echo "❌ Coverage $$coverage% is below 80%"; \
	else \
		echo "✅ Coverage $$coverage% meets threshold"; \
	fi