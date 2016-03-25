PRODUCT := strum
VERSION := $(shell cat VERSION)
WEBSITE := acksin.com

all: build

build: deps test
	go build -ldflags "-X main.version=$(VERSION)"
	$(MAKE) website-assets

deps:
	go get -u github.com/golang/lint/golint
	go get ./...

dev-deps:
	sudo apt-get install -y inkscape 

test:
	golint ./...
	go test -cover ./...
	go tool vet **/*.go

archive:
	tar cvzf $(PRODUCT)-$(VERSION).tar.gz $(PRODUCT)

release: spell build archive
	sed -i -e "s/^VERSION=.*$\/VERSION=$(VERSION)/g" website/install.sh
	git add website/install.sh
	-git commit -m "Version $(VERSION)"
	-git tag v$(VERSION) && git push --tags
	s3cmd put --acl-public $(PRODUCT)-$(VERSION).tar.gz s3://assets.acksin.com/$(PRODUCT)/${VERSION}/$(PRODUCT)-$(shell uname)-$(shell uname -i)-${VERSION}.tar.gz

website-assets:
	emacs DOCUMENTATION.org --batch --eval '(org-html-export-to-html nil nil nil t)'  --kill
	echo "---" > website/docs.html.erb
	echo "title: Acksin STRUM Docs" >> website/docs.html.erb
	echo "layout: docs" >> website/docs.html.erb
	echo "---" >> website/docs.html.erb
	cat DOCUMENTATION.html >> website/docs.html.erb
	rm DOCUMENTATION.html
	# cd website && go run logo.go > logo.svg && inkscape -z -d 150 -e autotune.png logo.svg

website:
	echo "Nothin here govn'r"

spell:
	for i in $(shell ls website/*.html); do \
		aspell check --mode=html $$i; \
	done

.PHONY: website website-dev
