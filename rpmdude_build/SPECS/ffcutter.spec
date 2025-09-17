Name:           ffcutter
Version:        1.0.0
Release:        1%{?dist}
Summary:        ffcutter system tray application
License:        GPLv3

Source0:        ffcutter
Source1:        ffcutter.png
Source2:        ffcutter.desktop

%description
Video cutter that works with ffmpeg.

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin/
mkdir -p %{buildroot}/usr/share/icons/
mkdir -p %{buildroot}/usr/share/applications/
install -m 755 %{SOURCE0} %{buildroot}/usr/bin/
install -m 644 %{SOURCE1} %{buildroot}/usr/share/icons/
install -m 644 %{SOURCE2} %{buildroot}/usr/share/applications/

%files
/usr/bin/ffcutter
/usr/share/icons/ffcutter.png
/usr/share/applications/ffcutter.desktop

%changelog
* Mon Jun 02 2025 Developer <github.com/sunaipa5> - 1.0-1
- Initial RPM release
