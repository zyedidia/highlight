package highlight

import "bytes"

func DetectFiletype(defs []*Def, filename string, fileSrc []byte) *Def {
	firstLine := bytes.Split(fileSrc, []byte("\n"))[0]

	for _, d := range defs {
		if d.ftdetect[0].Match([]byte(filename)) {
			return d
		}
		if len(d.ftdetect) > 1 {
			if d.ftdetect[1].Match(firstLine) {
				return d
			}
		}
	}

	return nil
}
