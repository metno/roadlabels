#!/bin/bash
echo "Docker.entry.sh"
set -e

#if [ ! -e  /var/lib/roadlabels/roadcams.db ]
then
    cp -v /roadlabels/roadcams.db /var/lib/roadlabels/roadcams.db
else 
    echo /var/lib/roadlabels/roadcams.db exist. Skip copy 
fi

if [ ! -e  /var/lib/roadlabels/users.db ]
then
    cp -v /roadlabels/userdb-empty.db /var/lib/roadlabels/users.db
else 
    echo /var/lib/roadlabels/users.db. Skip copy 
fi

/app/roadlabels
