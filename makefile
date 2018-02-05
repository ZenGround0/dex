gx_version=v0.12.1
gx-go_version=v1.6.0

deptools=deptools

gx=gx_$(gx_version)
gx-go=gx-go_$(gx-go_version)
gx_bin=$(deptools)/$(gx)
gx-go_bin=$(deptools)/$(gx-go)
bin_env=$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH)

all: dex

dex: deps
	go build

$(gx_bin):
	@echo "Downloading gx"
	@mkdir -p ./$(deptools)
	@rm -f $(deptools)/gx
	@wget -nc -q -O $(gx_bin).tgz https://dist.ipfs.io/gx/$(gx_version)/$(gx)_$(bin_env).tar.gz
	@tar -zxf $(gx_bin).tgz -C $(deptools) --strip-components=1 gx/gx
	@mv $(deptools)/gx $(gx_bin)
	@ln -s $(gx) $(deptools)/gx
	@rm $(gx_bin).tgz

$(gx-go_bin):
	@echo "Downloading gx-go"
	@mkdir -p ./$(deptools)
	@rm -f $(deptools)/gx-go
	@wget -nc -q -O $(gx-go_bin).tgz https://dist.ipfs.io/gx-go/$(gx-go_version)/$(gx-go)_$(bin_env).tar.gz
	@tar -zxf $(gx-go_bin).tgz -C $(deptools) --strip-components=1 gx-go/gx-go
	@mv $(deptools)/gx-go $(gx-go_bin)
	@ln -s $(gx-go) $(deptools)/gx-go
	@rm $(gx-go_bin).tgz


gx: $(gx_bin) $(gx-go_bin)

deps: gx
	$(gx_bin) install --global
	$(gx-go_bin) --verbose rewrite

check:
	go vet ./...
	golint -set_exit_status -min_confidence 0.3 ./...

test: deps
	go test -v ./...

clean: gx
	$(gx-go_bin) rewrite --undo
	@rm -rf $(deptools)

.PHONY: gx deps test 
