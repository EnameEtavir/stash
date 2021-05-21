package autotag

import (
	"github.com/stashapp/stash/pkg/models"
	"strings"
)

// QualifiedName has the Qualified meta information if the name String is judged to be good for robust matching
type QualifiedName struct {
	String    string
	Qualified bool
}

// Returns you the name and a qualification assessment if the name is good for "robust" matching
//
// This initial implementation is based just a rough guess on data I've seen so far of matches against
// FreeOnes names/aliases for performers where there are many one word aliases and sometimes names.
// FreeOnes-Babe "Belinda" is one example of a performer that is more widely known
// by here aliases that are also more "robust" , in contrast the name "Belinda" which is not judged to be robust
// by this current implementation, would produce a lot false-positives.
//
// This initial change to auto-tagging will only use qualified/robust names
// but later implementations should add "possible" matches with for example getPossiblePerformerMatches
// for performers_possible_scenes, performers_possible_images ...
func getQualifiedPerformerName(name string) QualifiedName {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "()-[]<>+")
	words := strings.Fields(name)
	isGood := len(words) > 1 && len(words[0]) > 1 && len(words[1]) > 1
	return QualifiedName{
		String:    name,
		Qualified: isGood,
	}
}

// Returns you the performers aliases and a qualification assessment
// for each if the name is good for "robust" matching
//
// This initial implementation is based just a rough guess on data I've seen so far of matches against
// FreeOnes names/aliases for performers where there are many one word aliases and sometimes names.
// FreeOnes-Babe "Belinda" is one example of a performer that is more widely known
// by here aliases that are also more "robust", in contrast the name "Belinda" which is not judged to be robust
// by this current implementation, would produce a lot false-positives.
//
// This initial change to auto-tagging will only use qualified/robust names
// but later implementations should add "possible" matches with for example getPossiblePerformerMatches
// for performers_possible_scenes, performers_possible_images ...
func getQualifiedPerformerAliases(performer *models.Performer) []QualifiedName {
	var ret []QualifiedName

	aliasesStr := strings.TrimSpace(performer.Aliases.String)
	aliases := strings.Split(aliasesStr, ",")
	for _, aliasStr := range aliases {
		aliasStr = strings.TrimSpace(aliasStr)
		if len(aliasStr) == 0 {
			continue
		}
		ret = append(ret, getQualifiedPerformerName(aliasStr))
	}

	return ret
}

// Returns you the performers names (name + aliases) and a qualification assessment
// for each if the name is good for robust matching
//
// This initial implementation is based just a rough guess on data I've seen so far of matches against
// FreeOnes names/aliases for performers where there are many one word aliases and sometimes names.
// FreeOnes-Babe "Belinda" is one example of a performer that is more widely known
// by here aliases that are also more "robust", in contrast the name "Belinda" which is not judged to be robust
// by this current implementation, would produce a lot false-positives.
//
// This initial change to auto-tagging will only use qualified/robust names
// but later implementations should add "possible" matches with for example getPossiblePerformerMatches
// for performers_possible_scenes, performers_possible_images ...
func getQualifiedPerformerNames(performer *models.Performer) []QualifiedName {
	var ret []QualifiedName
	ret = append(ret, getQualifiedPerformerName(performer.Name.String))
	ret = append(ret, getQualifiedPerformerAliases(performer)...)

	return ret
}
