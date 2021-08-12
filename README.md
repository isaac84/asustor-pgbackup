# asustor-pgbackup
Simple golang script to get a list of absolute file paths for all images and videos in the Asustor Photo Gallery  with a given tag applied to them.

*Usage:*

```./pgbackup /volume1/.@plugins/AppCentral/photogallery/etc/photogallery.db tagname```

The above creates a copy of the PhotoGallery sqllite3 database (as its locked) and finds all the images and photos with the specifited tag applied to them.
It then outputs each absolute file path to `backupfiles.txt` in the current directory.

This list can then be used to with `rsync --files-from=/tmp/pgbackup/backupfiles.txt` to copy the files to another directory setup with DataSync Center (eg. Google Drive)

The above can then be automated using bash and cron if required.

*Compiling:*

Compiling this for the Asustor NAS (AS5304T with ADM 3.5) can be tricky well at least on MacOS it was for me. So I used docker to cross compile it.

``` docker run --rm -v $PWD:/full/path/to/pgbackup/ -w /full/path/to/pgbackup/ golang:1.16.7-stretch go build -v```
