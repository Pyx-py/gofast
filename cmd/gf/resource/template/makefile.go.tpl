RPMROOT:=~/rpmbuild
PROJECT:={{.ProjectName}}
VERSION:=0.0.0.1

define RPMBUILD
	@rm -rf $(RPMROOT)/SOURCES/$(1)*
	@mkdir -p $(1)-$(2)
	@cp -f $(1) $(1)-$(2)
	@cp -f conf/$(1).conf $(1)-$(2)
	@cp -f centos7/$(1).service $(1)-$(2)
	@tar czvf $(1)-$(2).tar.gz $(1)-$(2)
	@rm -rf $(1)-$(2)
	@mv $(1)-$(2).tar.gz $(RPMROOT)/SOURCES
	@cp -f centos7/$(1).spec $(RPMROOT)/SOURCES
	@sed -i "s/PROJECT/$(1)/g" $(RPMROOT)/SOURCES/$(1).spec
	@sed -i "s/MVERSION/$(2)/g" $(RPMROOT)/SOURCES/$(1).spec
	rpmbuild -bb $(RPMROOT)/SOURCES/$(1).spec
	@rm -rf $(1)
	@mkdir -p rpms
	@mv $(RPMROOT)/RPMS/x86_64/$(1)-$(2)-* ./rpms
endef

all:
	go build -o $(PROJECT) main.go

clean:
	@rm -rf $(PROJECT) rpms

rpm:
	go build -o $(PROJECT) main.go
	@mkdir -p $(RPMROOT)/BUILD
	@mkdir -p $(RPMROOT)/BUILDROOT
	@mkdir -p $(RPMROOT)/RPMS
	@mkdir -p $(RPMROOT)/SOURCES
	@mkdir -p $(RPMROOT)/SPECS
	$(call RPMBUILD,$(PROJECT),$(VERSION))

prepare:
	find . -name "*.go"|xargs sed -i "s/GF_PROJECT_NAME/$(PROJECT)/g"
	mv -f conf/GF_PROJECT_NAME.conf conf/$(PROJECT).conf
	mv -f centos7/GF_PROJECT_NAME.service centos7/$(PROJECT).service
	mv -f centos7/GF_PROJECT_NAME.spec centos7/$(PROJECT).spec
	sed -i "s/GF_PROJECT_NAME/$(PROJECT)/g" centos7/$(PROJECT).service
	# mv ../$(shell basename $(shell pwd)) ../$(PROJECT)
	go mod init {{.ModuleName}}
	go mod tidy

usage:
	@echo "Usage:"
	@echo "                 build binary: make"
	@echo "                 build clean: make clean"
	@echo "                 build rpm: make rpm VERSION=x.x.x.x"

