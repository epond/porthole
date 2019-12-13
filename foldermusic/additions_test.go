package foldermusic

import (
	"testing"

	"github.com/epond/porthole/shared"
	"github.com/epond/porthole/test"
)

func TestNoMissingAlbums(t *testing.T) {
	scanned := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Shake - Iconoclastic Diaries"},
	}
	known := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Shake - Iconoclastic Diaries"},
	}
	missing := findMissingAlbums(scanned, known)
	test.ExpectInt(t, "number of missing albums", 0, len(missing))
}

func TestOneMissingAlbum(t *testing.T) {
	scanned := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Shake - Iconoclastic Diaries"},
	}
	known := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"The Krankies - It's Fan-dabi-dozi!"},
		shared.Album{"Shake - Iconoclastic Diaries"},
	}
	missing := findMissingAlbums(scanned, known)
	test.ExpectInt(t, "number of missing albums", 1, len(missing))
	test.Expect(t, "missing album", "The Krankies - It's Fan-dabi-dozi!", missing[0].Text)
}

func TestTwoMissingAlbums(t *testing.T) {
	scanned := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"Shake - Iconoclastic Diaries"},
	}
	known := []shared.Album{
		shared.Album{"Daniel Menche - Vent"},
		shared.Album{"The Krankies - It's Fan-dabi-dozi!"},
		shared.Album{"Shake - Iconoclastic Diaries"},
		shared.Album{"Throbbing Gristle - Discipline"},
	}
	missing := findMissingAlbums(scanned, known)
	test.ExpectInt(t, "number of missing albums", 2, len(missing))
	test.Expect(t, "missing album 1", "The Krankies - It's Fan-dabi-dozi!", missing[0].Text)
	test.Expect(t, "missing album 2", "Throbbing Gristle - Discipline", missing[1].Text)
}
