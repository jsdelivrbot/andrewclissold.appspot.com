#!/bin/bash

rsync -arv               \
    --exclude deploy.sh  \
    --exclude .git       \
    --exclude .gitignore \
    --exclude README.md  \
    --exclude *.swp      \
    --delete             \
    . ajclisso@login.secs.oakland.edu:public_html
