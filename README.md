# livego
## Simple and efficient live server:
* Installation and use is very simple;
* Writin in pure Golang, high performance, cross-platform;
* Support common transmission protocols, file formats, and encoding formats;

#### Supported transport protocols
- RTMP
- AMF
- HLS
- HTTP-FLV

#### Supported container formats
- FLV
- TS

#### Supported encoding formats
- H264
- AAC
- MP3

## Installation
After downloading the compiled binary directly, execute it on the command line.(https://github.com/bcrowe/livego/releases)

#### Download Source Code git clone
1. Go to the livego directory and execute go build `git clone https://github.com/gwuhaolin/livego.git`

## Getting Started
- Start the service: execute the livego binary to start the livego service;
- Upstream Push: Push the video stream to rtmp://localhost:1935/live/movie via RTMP, for example using ffmpeg -re -i demo.flv -c copy -f flv rtmp://localhost:1935/live/ Movie push
- Downstream playback: The following three playback protocols are supported. The playback address is as follows:
    * RTMP: rtmp://localhost:1935/live/movie
    * FLV: http://127.0.0.1:7001/live/movie.flv
    *   HLS: http://127.0.0.1:7002/live/movie.m3u8


### [Use with flv.js](https://github.com/gwuhaolin/blog/issues/3)

