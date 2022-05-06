package utils

import (
	"io"
	"os"
)

func Tail(file string, size int) (buf []byte, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	fs := fi.Size()
	if fs <= 0 {
		return
	}

	ss := fs
	bs := int64(4096)

	fb := make([]byte, bs)
	fp := fs
	ok := true
	for ok {
		if fp <= bs {
			fp = 0
			ok = false
		} else {
			fp -= bs
		}

		_, err = f.Seek(fp, io.SeekStart)
		if err != nil {
			return
		}

		n, e := f.Read(fb)
		if e != nil {
			return nil, e
		}
		if n > int(ss) {
			n = int(ss)
		}
		for i := n - 1; i >= 0; i-- {
			if fb[i] == '\n' {
				if ss != fs {
					size--
					if size < 1 {
						ok = false
						break
					}
				}
			}

			ss--
		}
	}

	if fs > ss {
		buf = make([]byte, fs-ss)
		_, err = f.ReadAt(buf, ss)
	}

	return
}
