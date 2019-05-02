package algo

import (
	"infra/appengine/rotang"
	"sort"
	"time"
)

// TZRecent implements a rota Generator scheduling members according to their timezones.
type TZRecent struct {
}

var _ rotang.RotaGenerator = &TZRecent{}

// NewTZRecent returns and instance of the TZRecent generator.
func NewTZRecent() *TZRecent {
	return &TZRecent{}
}

// Generate generates rotations fairly per time-zone.
func (t *TZRecent) Generate(sc *rotang.Configuration, start time.Time, previous []rotang.ShiftEntry, members []rotang.Member, shiftsToSchedule int) ([]rotang.ShiftEntry, error) {
	// Turns []rotang.Member into [][]rotang.Member with a slice of members per TZ.
	// The reason for not using a map here is to have the order of TZs consistent.
	sort.Slice(members, func(i, j int) bool {
		return members[i].TZ.String() < members[j].TZ.String()
	})
	var tzMembers [][]rotang.Member
	lastSeen := ""
	for _, m := range members {
		if lastSeen != m.TZ.String() {
			tzMembers = append(tzMembers, []rotang.Member{m})
			lastSeen = m.TZ.String()
			continue
		}
		tzMembers[len(tzMembers)-1] = append(tzMembers[len(tzMembers)-1], m)
	}

	// Since a pointer is used for Generate implying Generate won't change it up; better copy it.
	scCopy := *sc
	scCopy.Config.Shifts.ShiftMembers = 1

	recentGen := NewRecent()

	var perTZShifts [][]rotang.ShiftEntry
	for _, ms := range tzMembers {
		shifts, err := recentGen.Generate(&scCopy, start, previous, ms, shiftsToSchedule)
		if err != nil {
			return nil, err
		}
		perTZShifts = append(perTZShifts, shifts)
	}
	return tzSlice(perTZShifts)
}

// Name returns the name of the Generator.
func (t *TZRecent) Name() string {
	return "TZRecent"
}
