KENLM_DIR := third_party/kenlm
KENLM_BUILD := $(KENLM_DIR)/build
KENLM_MAX_ORDER ?= 6
PATCHES := $(sort $(wildcard patches/*.patch))

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
  BOOST_PREFIX ?= $(shell brew --prefix boost 2>/dev/null)
  CMAKE_EXTRA := -DBOOST_ROOT=$(BOOST_PREFIX) -DCMAKE_OSX_DEPLOYMENT_TARGET=12.0
else
  CMAKE_EXTRA :=
endif

.PHONY: all build-kenlm patch-kenlm test clean distclean dist-libs

all: build-kenlm

# Apply each patch if not already applied (idempotent via `git apply -R --check`).
patch-kenlm:
	@if [ ! -f $(KENLM_DIR)/CMakeLists.txt ]; then \
	  echo "kenlm submodule missing; run: git submodule update --init"; exit 1; \
	fi
	@for p in $(PATCHES); do \
	  abs=$$(cd $$(dirname $$p) && pwd)/$$(basename $$p); \
	  if (cd $(KENLM_DIR) && git apply -R --check "$$abs" >/dev/null 2>&1); then \
	    echo "  [skip] $$p (already applied)"; \
	  else \
	    echo "  [apply] $$p"; \
	    (cd $(KENLM_DIR) && git apply "$$abs"); \
	  fi; \
	done

build-kenlm: $(KENLM_BUILD)/lib/libkenlm.a

$(KENLM_BUILD)/lib/libkenlm.a: patch-kenlm
	cmake -S $(KENLM_DIR) -B $(KENLM_BUILD) \
	  -DCMAKE_BUILD_TYPE=Release \
	  -DFORCE_STATIC=ON \
	  -DKENLM_MAX_ORDER=$(KENLM_MAX_ORDER) \
	  -DENABLE_INTERPOLATE=OFF \
	  $(CMAKE_EXTRA)
	cmake --build $(KENLM_BUILD) --target kenlm kenlm_util -j

test: build-kenlm
	go test -count=1 ./...

# Package prebuilt static libs + headers for release.
# Output: dist/go-kenlm-libs-<os>-<arch>.tar.gz
DIST_OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
DIST_ARCH ?= $(shell uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')
DIST_NAME := go-kenlm-libs-$(DIST_OS)-$(DIST_ARCH)

dist-libs: build-kenlm
	rm -rf dist/$(DIST_NAME) dist/$(DIST_NAME).tar.gz
	mkdir -p dist/$(DIST_NAME)/lib dist/$(DIST_NAME)/include
	cp $(KENLM_BUILD)/lib/libkenlm.a $(KENLM_BUILD)/lib/libkenlm_util.a dist/$(DIST_NAME)/lib/
	cp -r $(KENLM_DIR)/lm $(KENLM_DIR)/util dist/$(DIST_NAME)/include/
	find dist/$(DIST_NAME)/include -name '*.cc' -delete
	tar -C dist -czf dist/$(DIST_NAME).tar.gz $(DIST_NAME)
	@echo "wrote dist/$(DIST_NAME).tar.gz"

clean:
	rm -rf $(KENLM_BUILD) dist

distclean: clean
	go clean -cache -testcache 2>/dev/null || true
