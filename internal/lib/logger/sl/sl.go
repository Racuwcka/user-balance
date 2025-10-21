package sl

func Err(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
