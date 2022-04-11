.PHONY: build clean test

hash = $(shell git log --pretty=format:'%h' -n 1)

build: clean
	mkdir -p build
	go build -o build/narrative -ldflags "-X main.build=${hash}" cmd/narrative/main.go

test:
	go test ./...

run:
	go run cmd/narrative/main.go

clean: clean_run
	rm -rf build

##################################################################

PROJECT_EXT=.tf
PROJECT_DIR=../aws-infrastructure/
OUTPUT=final

STAGE0=stage0
STAGE1=stage1

stage0:
	cat `find $(PROJECT_DIR)                  \
		-name "*$(PROJECT_EXT)"                 \
		-not -path "$(PROJECT_DIR).terraform/*" \
		| sort` > $(STAGE0)$(PROJECT_EXT)

stage1:
	./build/narrative   \
		-b "/*"           \
		-e "*/"           \
		-o $(STAGE1).md   \
		-i ./$(STAGE0)$(PROJECT_EXT)

stage2:
# sudo apt-get install groff pandoc
# pandoc -t ms -f markdown out.md -o temp.pdf
# --bibliography testlib.bib
	pandoc -s -t pdf -f markdown $(STAGE1).md -o $(OUTPUT).pdf
#	pandoc -s -t pdf -f markdown+raw_tex out.md -o temp.pdf

html:
	pandoc -s -t html $(STAGE1).md -o $(OUTPUT).html

man:
	pandoc -s -t man -f markdown $(STAGE1).md -o $(OUTPUT).1
	gzip $(OUTPUT).1
#	man ./$(OUTPUT).1.gz

slides:
	pandoc -s -t beamer $(STAGE1).md -o $(OUTPUT).pdf

clean_run:
	rm -f $(STAGE0)$(PROJECT_EXT)
	rm -f $(STAGE1).md
	rm -f $(OUTPUT).*
