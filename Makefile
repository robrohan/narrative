.PHONY: build clean test

hash = $(shell git log --pretty=format:'%h' -n 1)
PANDOC?=pandoc

OUTPUT_NAME?=manual
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

clean:
	rm -rf build
	rm -f $(PROJECT_DIR)/$(OUTPUT_NAME).*

docker_build:
	docker build . -t narrative

install_linux:
	apt-get install groff pandoc texlive-xetex

####################################################################
# Examples of creating outputs using local installed pandoc

to_pdf: run
# We cd here so that we can include bibliography files using
# an include path in the header that make sense
	cd $(PROJECT_DIR); \
	$(PANDOC) --pdf-engine=xelatex -s -t pdf \
		--citeproc \
		-f markdown $(OUTPUT_NAME).md \
		-o $(OUTPUT_NAME).pdf

to_manpage: run
	cd $(PROJECT_DIR); \
	$(PANDOC) -s -t man \
		--citeproc \
		-f markdown $(OUTPUT_NAME).md \
		-o $(OUTPUT_NAME).1 \
	; \
	gzip $(OUTPUT_NAME).1
# Example reading the output:
#	man $(PROJECT_DIR)/$(OUTPUT_NAME).1.gz

to_html: run
	cd $(PROJECT_DIR); \
	$(PANDOC) -s -t html \
		--citeproc \
		-f markdown $(OUTPUT_NAME).md \
		-o $(OUTPUT_NAME).html

####################################################################
# Examples of creating outputs using docker containers

docker_run:
	rm -f ./testdata/manual.md
	docker run --rm -it \
		-v $(shell pwd):/root/workspace \
		narrative \
			-i ./workspace/testdata/NARRATIVE \
			-o ./workspace/testdata/manual.md

docker_run_pdf: docker_run
	docker run --rm -it \
		-v $(shell pwd)/testdata:/root/workspace \
		robrohan/pandoc --pdf-engine=xelatex -s -t pdf \
		--citeproc \
		-f markdown ./manual.md \
		-o ./manual.pdf
