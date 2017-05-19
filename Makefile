
#
# Build the Alice Looking Glass
# -----------------------------
#
#

backend_dev:
	$(MAKE) -C backend/
	mv backend/alice-lg-* bin/


