.PHONY: analyze test-unit test-e2e test-goldens-update test-all clean

## Run flutter analyze
analyze:
	flutter analyze

## Run unit + widget tests
test-unit:
	flutter test

## Run E2E integration tests (requires connected device/emulator)
test-e2e:
	flutter test integration_test/

## Update golden files (run this to regenerate baselines)
test-goldens-update:
	flutter test --update-goldens

## Run all checks: analyze + unit tests
test-all: analyze test-unit

## Get dependencies
deps:
	flutter pub get

## Clean build artifacts
clean:
	flutter clean
	flutter pub get
