set -e

URI_C=https://raw.githubusercontent.com/cr-marcstevens/sha1collisiondetection/b45fcefc71270d9a159028c22e6d36c3817da188/lib/ubc_check.c
URI_H=https://raw.githubusercontent.com/cr-marcstevens/sha1collisiondetection/b45fcefc71270d9a159028c22e6d36c3817da188/lib/ubc_check.h
curl -s ${URI_C} -o ubc_check.c
curl -s ${URI_H} -o ubc_check.h
sha256sum --check sum.txt

cc ubc_check.c -o ubc_check.o -c
c++ -o verify.x verify.cpp ubc_check.o
./verify.x
