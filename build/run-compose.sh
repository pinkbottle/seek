#!/bin/bash
echo "🌘building crawler"
#run bash script and get its output
SHA=$(bash ./build-crawler.sh)
echo "🌗setting CRAWLER_IMAGE"
export CRAWLER_IMAGE=$SHA
echo "🌕running docker-compose.yaml"
cd ../ && docker-compose -f docker-compose.yaml up -d