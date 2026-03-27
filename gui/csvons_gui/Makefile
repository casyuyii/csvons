.PHONY: bootstrap clean preflight deps analyze format test check

bootstrap:
	./tool/bootstrap.sh

clean:
	./tool/clean_generated.sh

preflight:
	@command -v flutter >/dev/null 2>&1 || (echo "error: flutter is not installed or not on PATH" && exit 1)
	@command -v dart >/dev/null 2>&1 || (echo "error: dart is not installed or not on PATH" && exit 1)

deps: preflight
	flutter pub get

analyze: preflight
	flutter analyze

format: preflight
	dart format --set-exit-if-changed lib test

test: preflight
	flutter test

check: deps analyze format test
