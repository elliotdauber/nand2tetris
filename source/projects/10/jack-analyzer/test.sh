filepath=$1

go build
./jack-analyzer "$filepath.jack"
../../../tools/TextComparer.sh "$filepath.xml" "$filepath.X.xml"
rm -f "$filepath.X.xml"
