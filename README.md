# check_stable

## Why?
Because i have a customer with extremely heavily loaded machines.
There are checks based on check_logfiles which scan logfiles of
nightly batches. Sometimes the batches use so much cpu and memory
that check_by_ssh exits with 
Return code of ... is out of bounds

## How to configure it

## Code layout

    $GOPATH/
        src/
            github.com/
                lausser/
                    check_stable/
                        .git/
                        check_stable.go
                        check_stable_test.go
                        README.md
    
## Setup the Workspace

    mkdir ~/myworkspace
    export GOPATH=~/myworkspace
    cd $GOPATH
    mkdir -p src/github.com/lausser
    cd src/github.com/lausser
    git clone https://github.com/lausser/check_stable.git

## Build

    cd $GOPATH/src/github.com/lausser/check_stable
    go install
    $GOPATH/bin/check_stable

