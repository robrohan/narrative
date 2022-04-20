.PHONY: build clean test

hash = $(shell git log --pretty=format:'%h' -n 1)
PANDOC?=pandoc

OUTPUT_NAME?=final
PROJECT_DIR?=./testdata

build: clean
	mkdir -p build
	go build -o build/narrative -ldflags "-X main.build=${hash}" cmd/narrative/main.go

test:
	go test ./...

# Combine all the source files into a single file
# that can be used by all the to_* tasks below
run:
	rm -f $(PROJECT_DIR)/$(OUTPUT_NAME).md
	go run cmd/narrative/main.go \
		-i $(PROJECT_DIR)/NARRATIVE \
		-o $(PROJECT_DIR)/$(OUTPUT_NAME).md

to_pdf:
# We cd here so that we can include bibliography files using
# an include path in the header that make sense
	cd $(PROJECT_DIR); \
	$(PANDOC) --pdf-engine=xelatex -s -t pdf \
		--citeproc \
		-f markdown $(OUTPUT_NAME).md \
		-o $(OUTPUT_NAME).pdf

to_manpage:
	$(PANDOC) -s -t man \
		-f markdown $(PROJECT_DIR)/$(OUTPUT_NAME).md \
		-o $(PROJECT_DIR)/$(OUTPUT_NAME).1
	gzip $(PROJECT_DIR)/$(OUTPUT_NAME).1
# Example reading:
#	man $(PROJECT_DIR)/$(OUTPUT_NAME).1.gz

to_html:
	$(PANDOC) -s -t html \
		-f markdown $(PROJECT_DIR)/$(OUTPUT_NAME).md \
		-o $(PROJECT_DIR)/$(OUTPUT).html

clean:
	rm -rf build
	rm -f ./testdata/final.*

install_linux:
	apt-get install groff pandoc texlive-xetex
