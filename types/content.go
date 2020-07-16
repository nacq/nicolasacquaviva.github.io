package types

type Content struct {
	ParentDir string `bson:"parentDir"`
	Name      string `bson:"name"`
	Type      string `bson:"type"`
	Content   string `bson:"content"`
	Path      string `bson:"path"`
}
