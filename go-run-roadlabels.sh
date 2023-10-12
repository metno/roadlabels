#!/bin/env bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

export PYTHONPATH=$SCRIPT_DIR/exttools:$PYTHONPATH
if [ -z "$S3SecretKey" ]
then
    echo "\$S3SecretKey is empty. Please set"
    exit 1
fi

if [ -z "$S3AccessKey" ]
then
    echo "\$S3AccessKey is empty. Please set"
    exit 1
fi

time go run roadlabels.go -db-path var/lib/roadlabels/roadcams.db -userdb-path var/lib/roadlabels/users.db -access-key $S3AccessKey  -secret-key $S3SecretKey
                                                                     
