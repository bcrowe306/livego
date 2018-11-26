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

## 使用
2. 启动服务：执行 `livego` 二进制文件启动 livego 服务；
3. 上行推流：通过 `RTMP` 协议把视频流推送到 `rtmp://localhost:1935/live/movie`，例如使用 `ffmpeg -re -i demo.flv -c copy -f flv rtmp://localhost:1935/live/movie` 推送；
4. 下行播放：支持以下三种播放协议，播放地址如下：
    - `RTMP`:`rtmp://localhost:1935/live/movie`
    - `FLV`:`http://127.0.0.1:7001/live/movie.flv`
    - `HLS`:`http://127.0.0.1:7002/live/movie.m3u8`


### [和 flv.js 搭配使用](https://github.com/gwuhaolin/blog/issues/3)

对Golang感兴趣？请看[Golang 中文学习资料汇总](http://go.wuhaolin.cn/)
