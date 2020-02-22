#!/bin/sh
gcloud compute instances create thunt \
    --image cos-stable-80-12739-68-0 \
    --image-project cos-cloud \
    --metadata-from-file google-container-manifest=containers.yaml \
    --tags http-server \
    --zone us-central1-a \
    --machine-type f1-micro