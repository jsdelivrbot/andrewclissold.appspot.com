#!/bin/bash

minify css/styles.css
minify css/syntax.css
rsync -arv               \
    --exclude deploy.sh  \
    --exclude .git       \
    --exclude .gitignore \
    --exclude README.md  \
    --exclude *.swp      \
    --delete             \
    . ajclisso@login.secs.oakland.edu:public_html
