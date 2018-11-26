#!/bin/bash 

function exit_on_error {
    if [ $1 != 0 ]; then 
        exit $1
    fi 
}


# Lib 
echo "Testing lib"
cd ./lib 
go test 
exit_on_error $?

### 
# Modules 
### 

echo "Testing modules/compare" 
cd ../modules/compare 
go test 
exit_on_error $? 

echo "Testing modules/dal" 
cd ../../modules/dal
go test 
exit_on_error $? 

echo "Testing modules/gen" 
cd ../../modules/gen
go test 
exit_on_error $? 

