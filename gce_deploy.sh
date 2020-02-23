#!/bin/sh
gcloud compute instances create-with-container thunt \
    --container-image=registry.hub.docker.com/lthrane/thunt