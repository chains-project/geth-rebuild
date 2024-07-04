# geth reference build

geth-reference: go1.22.3
Running `readelf -p .comment geth-reference`...

String dump of section '.comment':
[ 0] GCC: (Ubuntu 9.4.0-1ubuntu1~20.04.2) 9.4.0

# geth reproducing build

geth-reproduce: go1.22.3
Running `readelf -p .comment geth-reproduce`...

String dump of section '.comment':
[ 0] GCC: (Ubuntu 13.2.0-23ubuntu4) 13.2.0

Built in OS:
PRETTY_NAME="Ubuntu 24.04 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION="24.04 LTS (Noble Numbat)"
