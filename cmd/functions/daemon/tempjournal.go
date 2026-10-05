package daemon

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
)

// tempJournalFileProd records the allow rows of the live `edit-protected --forward` grant BEFORE
// they are written: pinned links keep enforcing after a SIGKILL, with those rows still in the
// pinned whitelist, so `daemon --lockdown` and the next start strip them (stripTempJournal). In the
// self-guarded /etc/app-listener for the same integrity reason as pin-state.json; a forged journal
// can only cause denials (StripPinnedTempAllows deletes, never adds).
const tempJournalFileProd = "/etc/app-listener/temp-grants.json"

var tempJournalFile = tempJournalFileProd

type tempJournalRow struct {
	Res uint32 `json:"res"`
	Dev uint64 `json:"dev"`
	Ino uint64 `json:"ino"`
}

// writeTempJournal records rows in place; nil empties the journal.
func writeTempJournal(rows []guard.GuardResInodeKey) error {
	if len(rows) == 0 {
		return writeInPlace(tempJournalFile, nil)
	}
	out := make([]tempJournalRow, len(rows))
	for i, r := range rows {
		out[i] = tempJournalRow{Res: r.ResId, Dev: r.Ino.Dev, Ino: r.Ino.Ino}
	}
	data, err := json.Marshal(out)
	if err != nil {
		return err
	}
	return writeInPlace(tempJournalFile, data)
}

func readTempJournal() ([]guard.GuardResInodeKey, error) {
	data, err := os.ReadFile(tempJournalFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	var in []tempJournalRow
	if err := json.Unmarshal(data, &in); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", tempJournalFile, err)
	}
	rows := make([]guard.GuardResInodeKey, len(in))
	for i, r := range in {
		rows[i] = guard.GuardResInodeKey{ResId: r.Res, Ino: guard.GuardInodeKey{Dev: r.Dev, Ino: r.Ino}}
	}
	return rows, nil
}

// stripTempJournal removes a dead daemon's journaled temporary allows from its pinned whitelist,
// then empties the journal. The journal is kept when nothing could be stripped (no trustworthy
// pin generation, or its owner still looks alive), so a later run retries.
func stripTempJournal(logPrefix, pinBase string) {
	rows, err := readTempJournal()
	if err != nil {
		log.Errorf("%s: CRITICAL: %v — temporary edit-protected allows a killed daemon left may still be "+
			"in its pinned whitelist until its pins are retired", logPrefix, err)
		return
	}
	if len(rows) == 0 {
		return
	}
	rec := recoverPinState(logPrefix)
	if rec.Gen == "" {
		return
	}
	if _, err := guard.StripPinnedTempAllows(guard.SharedPinPrefix(pinBase, rec.Gen), rows); err != nil {
		log.Errorf("%s: CRITICAL: stripping temporary edit-protected allows from generation %s: %v",
			logPrefix, rec.Gen, err)
		return
	}
	if err := writeTempJournal(nil); err != nil {
		log.Warnf("%s: clearing %s: %v", logPrefix, tempJournalFile, err)
	}
}
