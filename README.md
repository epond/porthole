# porthole

A dashboard presented as a web page that shows the most recent additions to a record collection. It is intended to run on a Raspberry Pi.

* Once porthole is installed on the Pi, every time the Pi is turned on it updates porthole to the latest version by pulling the source from GitHub and doing a full build, then starts the app.
* Golang was used because it is easily installed via apt-get and it is very fast both when compiling and running.
* A file specified by the KNOWN_RELEASES_FILE environment variable is used to keep track of what as been added to the record collection, in order. New additions are added to the bottom of the file.
* The record collection must be arranged in folders where each release is a subfolder of the artist responsible.
* The folder names are all that is used to determine the artist name and release name. The app does not care about audio files or embedded metadata.
* Use the FOLDERS_TO_SCAN environment variable to indicate any number of locations where releases are found, separated by commas.
* Each location in FOLDERS_TO_SCAN is split into `<root>:<depth>` where `<depth>` indicates how many levels beneath `<root>` the release folders can be found. This allows for splitting the collection by each artist's first letter.
