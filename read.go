package sumfs

func Read(dir string) (*FS, error) {
	fsys := NewFS()
	err := fsys.Read(dir)
	return fsys, err
}
