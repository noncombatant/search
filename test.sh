#!/bin/sh

function clean_up() {
  rm -rf test_files
}

function check() {
  test "$1" = "$2" || echo "Expected \"$1\" got \"$2\""
}

trap clean_up HUP INT QUIT ILL TRAP ABRT KILL

mkdir -p test_files/stuff test_files/things
touch test_files/hello.txt test_files/wow.pdf
touch test_files/things/noodles.jpg
echo "Flavors and tastes are a big part of food." > test_files/important
echo "Flavors and tastes are a not big part of food." > test_files/silly

result=$(./search -t d test_files)
expected="test_files
test_files/stuff
test_files/things"
check "$expected" "$result"

result=$(./search -t f test_files)
expected="test_files/hello.txt
test_files/important
test_files/silly
test_files/things/noodles.jpg
test_files/wow.pdf"
check "$expected" "$result"

result=$(./search -n pdf test_files)
expected="test_files/wow.pdf"
check "$expected" "$result"

result=$(./search -c flavor test_files)
expected="test_files/important:Flavors and tastes are a big part of food.
test_files/silly:Flavors and tastes are a not big part of food."
check "$expected" "$result"

result=$(./search -n '!silly' -c flavor test_files)
expected="test_files/important:Flavors and tastes are a big part of food."
check "$expected" "$result"

result=$(./search -n '\.jpg$' test_files)
expected="test_files/things/noodles.jpg"
check "$expected" "$result"

result=$(./search -s 5 test_files)
expected="test_files
test_files/important
test_files/silly
test_files/stuff
test_files/things"
check "$expected" "$result"

result=$(./search -s 5 -t f test_files)
expected="test_files/important
test_files/silly"
check "$expected" "$result"

result=$(./search -a 2021-10-24 test_files)
expected="test_files
test_files/hello.txt
test_files/important
test_files/silly
test_files/stuff
test_files/things
test_files/things/noodles.jpg
test_files/wow.pdf"
check "$expected" "$result"

result=$(./search -b 3031-10-24 test_files)
check "$expected" "$result"

result=$(./search -a 3031-10-24 test_files)
expected=""
check "$expected" "$result"

clean_up
