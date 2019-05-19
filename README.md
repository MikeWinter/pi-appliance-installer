# Raspberry Pi On-Boot Provisioning

[Installing Raspbian](https://www.raspberrypi.org/documentation/installation/installing-images/) is
as simple as writing the OS image to an SD card. But what about customising it? This gets tricky if
you aren’t already running a Linux OS.

This tool aims to make customising a Raspberry Pi easier for anyone using Windows or MacOS.

## How does this tool help?

The Raspbian OS image includes two partitions: root and boot. The root partition contains the OS
but it uses the ext4 disk format, rendering it inaccessible to some systems. We can’t customise it
if we can't write to it.

The boot partition uses the fat32 disk format, which is supported by almost all operating systems.
We use this to stage files we want included in the root partition, and copy them over before
Raspbian starts.

## How does it work?

Installing the tool modifies the normal [startup process](
https://en.wikipedia.org/wiki/Linux_startup_process) so that files can be copied from the boot
partition into the root partition. Once complete, startup is resumed.

## Prerequisites

* [Golang](https://golang.org/)
