#!/bin/bash
go get github.com/axw/gocov/...
go get github.com/AlekSi/gocov-xml

shopt -s nullglob

if [ ls *.coverprofile > /dev/null 2>&1 ];
then
    echo "No coverprofiles found. Run a test suite first."
    exit

else
    echo "Converting coverprofiles..."
    for file in *.coverprofile
    do
        echo "Converting $file from Go coverprofile to Cobertura XML format"
        XML_NAME=${file%.*}
        gocov convert $file | gocov-xml > $XML_NAME.xml
    done
fi
