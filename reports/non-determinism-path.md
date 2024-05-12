# geth reference build

Running `readelf -p .rodata geth-reference | grep /home/travis`...
[9e8b70] /home/travis/gopath/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/libusbi.h
[9e8be8] /home/travis/gopath/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/os/events_posix.c
[9e8c68] /home/travis/gopath/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/os/linux_netlink.c
[9e8ce8] /home/travis/gopath/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/os/linux_usbfs.c
[9e8d60] /home/travis/gopath/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/core.c
[9e8e00] /home/travis/gopath/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/hotplug.c
[9e8e78] /home/travis/gopath/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/io.c
[9e9db8] /home/travis/gopath/pkg/mod/github.com/ethereum/c-kzg-4844@v1.0.0/bindings/go/../../src/c_kzg_4844.c

# geth reproducing build

Running `readelf -p .rodata geth-reproduce | grep /root/go/pkg`...
[9f9db0] /root/go/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/libusbi.h
[9f9e18] /root/go/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/os/events_posix.c
[9f9e88] /root/go/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/os/linux_netlink.c
[9f9ef8] /root/go/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/os/linux_usbfs.c
[9f9f68] /root/go/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/core.c
[9fa000] /root/go/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/hotplug.c
[9fa068] /root/go/pkg/mod/github.com/karalabe/hid@v1.0.1-0.20240306101548-573246063e52/libusb/libusb/io.c
