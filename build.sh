#!/bin/bash

if [ ! -f "Dockerfile" ]; then
    curl -s https://gitlab.zalopay.vn/trungdv3repos/docker-templates/blob/master/dockerfile-golang-template?private_token=E7EWzo3-o4U9aCYs9k7s > Dockerfile
fi

for acc in $(printf "\n" | git credential-osxkeychain get); do 
    export $acc
done

if [[ "$username" == "" || "$password" == "" ]]; then
    echo "Error missing username or password";
    exit;
fi


docker build --build-arg=username --build-arg=password -t mas_image .