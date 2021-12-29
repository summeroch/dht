## Description

项目来源自[来源地址](https://github.com/shiyanhui/dht)

## Introduction

DHT实现了BitTorrent DHT协议，主要包括：

+ [BEP-3 (部分)](http://www.bittorrent.org/beps/bep_0003.html)
+ [BEP-5](http://www.bittorrent.org/beps/bep_0005.html)
+ [BEP-9](http://www.bittorrent.org/beps/bep_0009.html)
+ [BEP-10](http://www.bittorrent.org/beps/bep_0010.html)

它包含两种模式，标准模式和爬虫模式。标准模式遵循DHT协议，你可以把它当做一个标准的DHT组件。爬虫模式是为了嗅探到更多torrent文件信息，它在某些方面不遵循DHT协议。

## Installation

    go get github.com/summeroch/dht

## Usage

[releases](https://github.com/summeroch/dht-spider/releases) 页面下载对应的二进制文件(或自行编译)。将可执行文件与config.yml放置于同一路径下即可运行。