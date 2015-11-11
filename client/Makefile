PATH       := node_modules/.bin:$(PATH)
SHELL      := /usr/bin/env bash
# ------------------------------------------------------------------------------
bower      := $(shell cat .bowerrc | python -c 'import json,sys;print json.load(sys.stdin)["directory"];')

source     := src
source_js   = $(shell find $(source) ! -path "$(bower)/*" -type f -name '*.js')
source_css  = $(shell find $(source) ! -path "$(bower)/*" -type f -name '*.scss')
source_rest = $(shell find $(source) ! -path "$(bower)/*" ! -name '*.js' ! -name '*.scss' -type f)

output_dir := www
output_css  = $(output_dir)/css/index.css
output_js   = $(output_dir)/js/index.js
output_rest = $(patsubst $(source)/%,$(output_dir)/%,$(source_rest))

output_apk := ../build/android.apk

.PHONY: all
all: build

.PHONY: build
build:
	@echo "-----> Rebuilding client..."
	@$(MAKE) javascript stylesheet rest

.PHONY: emulate
emulate: node_modules platforms/android build
	@ionic run android

.PHONY: build.android
build.android: $(output_apk)

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
	rm -rf $(output_dir) $(output_apk)

node_modules: package.json
	npm install
	npm dedupe
	@touch $@

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

$(output_dir)/%: $(source)/%
	@mkdir -p $(dir $@)
	cp $< $@

$(output_apk): node_modules platforms/android build
	@mkdir -p $(dir $@)
	@echo 'Building APK...'
	@cordova build android
	@cp platforms/android/build/outputs/apk/android-armv7-debug.apk $(@)

platforms/%: config.xml node_modules
	@rm -rf $@
	@mkdir -p $(output) $(dir $@)
	@cordova platform add $*
	@ionic state restore
