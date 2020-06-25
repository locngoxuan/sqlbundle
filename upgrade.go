package sqlbundle

func (sb *SQLBundle) Upgrade() error {
	err := sb.readConfig()
	if err != nil {
		return err
	}

	script := &MigrationScript{
		Version:  sb.ReadVersion(),
		Group:    sb.Config.GroupId,
		Artifact: sb.Config.ArtifactId,
	}

	err = sb.collectMigrations(script)
	if err != nil {
		return err
	}
	//script.ListAll()
	return nil
}


