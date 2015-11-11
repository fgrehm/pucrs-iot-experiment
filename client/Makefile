PATH       := node_modules/.bin:$(PATH)
SHELL      := /usr/bin/env bash
# ------------------------------------------------------------------------------
bower      := $(shell cat .bowerrc | python -c 'import json,sys;print json.load(sys.stdin)["directory"];')

source     := www
source_js   = $(shell find $(source) ! -path "$(bower)/*" -type f -name '*.js')
source_css  = $(shell find $(source) ! -path "$(bower)/*" -type f -name '*.scss')
source_rest = $(shell find $(source) ! -path "$(bower)/*" ! -name '*.js' ! -name '*.scss' -type f)

output     := dist
output_css  = $(output)/css/index.css
output_js   = $(output)/js/index.js
output_rest = $(patsubst $(source)/%,$(output)/%,$(source_rest))

.PHONY: all
all: build

.PHONY: build
build:
	@echo "-----> Rebuilding client..."
	@$(MAKE) javascript stylesheet rest

.PHONY: javascript
javascript: $(output_js)

.PHONY: stylesheet
stylesheet: $(output_css)

.PHONY: rest
rest: $(output_rest)

.PHONY: install
install: node_modules $(bower)

.PHONY: clean
clean:
	rm -rf $(output)

node_modules: package.json
	npm install
	npm dedupe

$(bower): node_modules bower.json
	bower install --allow-root
	@touch $@

$(output_js): node_modules $(bower) $(source_js)
	@mkdir -p $(dir $@)
	@echo 'Executing browserify...'
	@browserify $(source)/index.js \
		--ignore-missing \
			--transform [ configurify --pattern '**/config/*' ] \
			--transform [ debowerify ] \
		--outfile $@

$(output_css): node_modules $(bower) $(source_css)
	@mkdir -p $(dir $@)
	@echo 'Executing node-sass...'
	@node-sass $(source)/index.scss | postcss \
		--use autoprefixer \
		--autoprefixer.browsers 'last 2 versions' \
		--output $@

$(output)/%: $(source)/%
	@mkdir -p $(dir $@)
	cp $< $@
