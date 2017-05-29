
#
# Build the Alice Looking Glass
# -----------------------------
#
#


all: alice

client_dev:
	$(MAKE) -C client/

client_prod:
	$(MAKE) -C client/ client_prod

backend_dev: client_dev
	$(MAKE) -C backend/


backend_prod: client_prod
	$(MAKE) -C backend/ bundle
	$(MAKE) -C backend/ linux
	

alice: client_prod backend_prod
	mv backend/alice-lg-* bin/

