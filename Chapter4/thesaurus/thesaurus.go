package thesaurus

// Thesaurus interface specifies a way of delivering
//  synonyms to program
type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}
