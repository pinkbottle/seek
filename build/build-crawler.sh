cd  ../cmd/crawler && SHA=$(KO_DOCKER_REPO=$DOCKER_USER ko publish . | tail -1)