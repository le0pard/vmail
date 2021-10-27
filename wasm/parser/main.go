package main

func ReportFromHTML(document []byte) (string, error) {
	parser, err := InitParser()
	if err != nil {
		return "", err
	}

	err = parser.Report(document)
	if err != nil {
		return "", err
	}

	return "", nil
}
