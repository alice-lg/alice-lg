#
# Build Helper
# ------------
#
# Create a full build by just typing make.
# This will automatically install all dependencies from NPM and
# start the build process.
#
# While developing, you might want to use 'make watch'
# which will automatically restart gulp in case something went
# wrong.
#

VERSION=$(shell cat ../VERSION)

DIST_BUILDS=../../birdseye-static/builds
DIST=birdseye-ui-dist-$(VERSION).tar.gz

# == END CONFIGURATION ==

DIST_BUILD=$(addprefix $(DIST_BUILDS)/, $(DIST))

all: deps client

deps:
	@echo "Installing dependencies"
	npm install

client:
	@echo "Building alice UI"
	gulp

client_prod:
	@echo "Building alice UI (production)"
	DISABLE_LOGGING=1 gulp

watch:
	while true; do gulp watch; done

$(DIST_BUILD): deps client_prod
	@echo "Creating alice ui distribution"
	tar cvzf $(DIST) build/
	mv $(DIST) $(DIST_BUILDS)
	@echo ""
	@echo "Done. Don't forget to push the dist to github"



dist: $(DIST_BUILD)
	echo $(DIST_BUILD)
