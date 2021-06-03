TAG = $(shell git tag | sort -r --version-sort | head -n1)
SEMVERS = $(shell echo $(TAG) | semver)

major:
	@set -e;\
	export major=$$(( $$(echo '$(SEMVERS)'|jq -r .major|sed 's/^v//') + 1 )); \
	git tag -m "v$$major.0.0 major release" --sign v$$major.0.0
minor:
	@set -e;\
	export major=$$(( $$(echo '$(SEMVERS)'|jq -r .major|sed 's/^v//') + 0 )); \
	export minor=$$(( $$(echo '$(SEMVERS)' | jq -r .majorminor | sed -e 's/^v[0-9]*\.//' -e 's/\.[0-9]$$//' ) + 1 ));\
	export patch=$$(( $$(echo '$(SEMVERS)' | jq -r .canonical|sed 's/^v[0-9]*\.[0-9]*\.//') + 0 ));\
	git tag -m "v$$major.$$minor.$$patch minor release" --sign v$$major.$$minor.$$patch
patch:
	@set -e;\
	export major=$$(( $$(echo '$(SEMVERS)'|jq -r .major|sed 's/^v//') + 0 )); \
	export minor=$$(( $$(echo '$(SEMVERS)' | jq -r .majorminor | sed -e 's/^v[0-9]*\.//' -e 's/\.[0-9]$$//' ) + 0 ));\
	export patch=$$(( $$(echo '$(SEMVERS)' | jq -r .canonical|sed 's/^v[0-9]*\.[0-9]*\.//') + 1 ));\
	git tag -m "v$$major.$$minor.$$patch patch fix" --sign v$$major.$$minor.$$patch

semver:
	git tag -f -m '$(TAG)' --sign $$(echo '$(SEMVERS)' | jq -r .major)
	git tag -f -m '$(TAG)' --sign $$(echo '$(SEMVERS)' | jq -r .majorminor)

release:
	git push
	git push --tags -f

install:
	go install ./cmd/...
