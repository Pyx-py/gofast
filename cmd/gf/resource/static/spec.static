%define debug_package %{nil}
Name: PROJECT
Summary: PROJECT
Version: MVERSION
License: MIT
Release: 1
Group: Applications/System

Source0: %{name}-%{version}.tar.gz

Requires(post): /usr/bin/systemctl
Requires(preun): /usr/bin/systemctl

%description
Include PROJECT

%preun
/usr/bin/systemctl stop PROJECT.service
/usr/bin/systemctl disable PROJECT.service

%prep
/usr/bin/mkdir -p /var/lib/PROJECT
%setup -q

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/etc/PROJECT
mkdir -p %{buildroot}/usr/lib/systemd/system
mkdir -p %{buildroot}/var/lib/PROJECT
mkdir -p %{buildroot}/var/log/PROJECT

cp $RPM_BUILD_DIR/${RPM_PACKAGE_NAME}-${RPM_PACKAGE_VERSION}/PROJECT $RPM_BUILD_ROOT/usr/bin
cp $RPM_BUILD_DIR/${RPM_PACKAGE_NAME}-${RPM_PACKAGE_VERSION}/PROJECT.service $RPM_BUILD_ROOT/usr/lib/systemd/system
cp $RPM_BUILD_DIR/${RPM_PACKAGE_NAME}-${RPM_PACKAGE_VERSION}/PROJECT.conf $RPM_BUILD_ROOT/etc/PROJECT

%files
%defattr(-, root, root, 0755)
/usr/bin/PROJECT
/usr/lib/systemd/system/PROJECT.service
%config(noreplace) /etc/PROJECT/PROJECT.conf
%changelog
