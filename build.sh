#!/bin/bash

rm -rf build

mkdir -p build/RGTV
mkdir -p build/Imgs

go build -o build/RGTV/RGTV
cp ffmpeg_bin/ffmpeg-arm64 build/RGTV/ffmpeg
cp sample_tv.json build/RGTV/tv.json

cp RGTV.sh build/
cp ./embeddata/RGTV.png build/Imgs

cd build
zip -r RGTV.zip *
