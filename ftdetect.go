package highlight

import "bytes"

func DetectFiletype(defs []*Def, filename string, fileSrc []byte) string {
	firstLine := bytes.Split(fileSrc, []byte("\n"))[0]

	for _, d := range defs {
		if d.ftdetect[0].Match([]byte(filename)) {
			return d.ft
		}
		if d.ftdetect[1].Match(firstLine) {
			return d.ft
		}
	}

	return "Unknown"
}
