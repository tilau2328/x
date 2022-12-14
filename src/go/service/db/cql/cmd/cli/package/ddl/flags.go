package ddl

type (
	KeySpaceFlags struct {
		Name        string
		Durable     bool
		Replication map[string]string
	}
	TableFlags struct {
		Name     string
		Keyspace string
		Fields   []string
	}
)

func ToKeySpace(f KeySpaceFlags) KeySpace {
	return KeySpace{}
}
func ToKeySpacePatch(KeySpaceFlags) []Patch {
	return nil
}
func ToTableKey(f TableFlags) TableKey {
	return TableKey{
		KeySpace: f.Keyspace,
		Name:     f.Name,
	}
}
func ToTable(f TableFlags) Table {
	return Table{}
}
func ToTablePatch(f TableFlags) []Patch {
	return nil
}
