## iTunes Export Console (golang)

A console application that exports iTunes Playlists using the iTunes Music Library.xml plist.

This is a port of the previous Scala version, found here: https://github.com/ericdaugherty/itunesexport-scala and .Net version, found here: http://www.ericdaugherty.com/dev/itunesexport/1.x/ 


## Compiling

```
go build 
```
## Usage

```
usage: itunesexport [<flags>] [include <playlist name>...]

Flags:
    -library <file path>        Path to iTunes Music Library XML File.
    -output <file path>         Path where the playlists should be written.
    -type <M3U|EXT|WPL|ZPL>     Type of playlist file to write.  Defaults to M3U
                                EXT = M3U Extended, WPL = Windows Playlist, ZPL = Zune Playlist
    -includeAll                 Include all user defined playlists.
    -includeAllWithBuiltin      Include All playlists, including iTunes defined playlists
    -copy <COPY TYPE>           Copy the music tracks as well, according the the COPY TYPE scheme...
        NONE                    (default) The music files will not be copied.                               
        PLAYLIST                Copies the music into a folder for each playlist.
        ITUNES                  Copies using the itunes music/<Artist>/<Album>/<Track> structure.
        FLAT                    Copies all the music into the output folder.
```
