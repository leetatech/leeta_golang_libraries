INSTALL_GOMARKDOC := go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

.PHONY: generate_readme
generate_readme:
	@command -v gomarkdoc > /dev/null || $(INSTALL_GOMARKDOC)
	gomarkdoc -o '{{ .Dir }}/README.md' ./...