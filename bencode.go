package byteslego

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type bencoder struct {
	bufio.Writer
}

func (b *bencoder) encode(d interface{}) (string, error) {
	var (
		out string
		err error
	)
	switch d.(type) {
	case int, int32, int64:
		out, err = b.encodeInt(int64(d.(int)))
	case string:
		out, err = b.encodeStr(d.(string))
	case []interface{}:
		out, err = b.encodeList(d.([]interface{}))
	case map[string]interface{}:
		out, err = b.encodeDict(d.(map[string]interface{}))
	}
	return out, err
}

func (b *bencoder) encodeInt(i int64) (string, error) {
	return fmt.Sprintf("i%de", i), nil
}

func (b *bencoder) encodeStr(s string) (string, error) {
	return fmt.Sprintf("%d:%s", len(s), s), nil
}

func (b *bencoder) encodeList(l []interface{}) (string, error) {
	var (
		s   []string
		err error
	)
	for _, e := range l {
		switch e.(type) {
		case int, int32, int64:
			i, err := b.encodeInt(int64(e.(int)))
			if err != nil {
				break
			}
			s = append(s, i)
		case string:
			es, err := b.encodeStr(e.(string))
			if err != nil {
				break
			}
			s = append(s, es)
		case []interface{}:
			el, err := b.encodeList(e.([]interface{}))
			if err != nil {
				break
			}
			s = append(s, el)
		case map[string]interface{}:
			ed, err := b.encodeDict(e.(map[string]interface{}))
			if err != nil {
				break
			}
			s = append(s, ed)
		}
	}

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("l%se", strings.Join(s, "")), nil
}

func (b *bencoder) encodeDict(d map[string]interface{}) (string, error) {
	var (
		s   []string
		err error
	)

	// sort keys first
	keys := make(sort.StringSlice, len(d))
	i := 0
	for k := range d {
		keys[i] = k
		i++
	}
	keys.Sort()

	for _, k := range keys {
		v := d[k]
		enckey, err := b.encodeStr(k)
		if err != nil {
			break
		}

		switch v.(type) {
		case string:
			es, err := b.encodeStr(v.(string))
			if err != nil {
				break
			}
			s = append(s, fmt.Sprintf("%s%s", enckey, es))
		case int, int32, int64:
			ei, err := b.encodeInt(int64(v.(int)))
			if err != nil {
				break
			}
			s = append(s, fmt.Sprintf("%s%s", enckey, ei))
		case []interface{}:
			el, err := b.encodeList(v.([]interface{}))
			if err != nil {
				break
			}
			s = append(s, fmt.Sprintf("%s%s", enckey, el))
		case map[string]interface{}:
			ed, err := b.encodeDict(v.(map[string]interface{}))
			if err != nil {
				break
			}
			s = append(s, fmt.Sprintf("%s%s", enckey, ed))
		case map[string][]interface{}:
			var ss []string
			val := v.(map[string][]interface{})
			for sk, sv := range val {
				ses, err := b.encodeStr(sk)
				if err != nil {
					break
				}
				sel, err := b.encodeList(sv)
				if err != nil {
					break
				}
				ss = append(ss, fmt.Sprintf("%s%s", ses, sel))
			}

			if err != nil {
				break
			}
			s = append(s, fmt.Sprintf("%sd%se", enckey, strings.Join(ss, "")))
		}
	}

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("d%se", strings.Join(s, "")), nil
}

type bdecoder struct {
	bufio.Reader
}

func (b *bdecoder) decode() (interface{}, error) {
	c, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	// we read first byte on start of the loop, so we unread it
	// and read from the start up-to delimiter since the size
	// can have more than 1 digit.
	if err = b.UnreadByte(); err != nil {
		return nil, err
	}

	var out interface{}
	switch c {
	case 'i':
		// then delegate to decoding integer
		out, err = b.decodeInt()
		if err != nil {
			break
		}
	case 'l':
		out, err = b.decodeList()
		if err != nil {
			break
		}

		// decoding a list puts it in a wrapper list,
		// so we extract that single item as the decoded list
		out = out.([]interface{})[0]
	case 'd':

	// case 'e': -- bencoding cannot start in e

	default:
		out, err = b.decodeStr()
		if err != nil {
			break
		}
	}

	return out, err
}

func (b *bdecoder) decodeStr() (string, error) {
	rsize, err := b.ReadString(':')
	if err != nil {
		return "", err
	}

	// since rsize will include colon, we strip it
	rsize = rsize[:len(rsize)-1]

	size, err := strconv.Atoi(rsize)
	if err != nil {
		return "", err
	}

	// read the next bytes based on the size
	var buf []byte
	for i := 0; i < size; i++ {
		rb, err := b.ReadByte()
		if err != nil {
			return "", err
		}
		buf = append(buf, rb)
	}

	return string(buf), nil
}

func (b *bdecoder) decodeInt() (int64, error) {
	// read int together with its delimeters
	rint, err := b.ReadString('e')
	if err != nil {
		return -1, err
	}

	// strip int delimeter and convert ti int64
	i, err := strconv.ParseInt(rint[1:len(rint)-1], 10, 64)
	if err != nil {
		return -1, err
	}

	return i, nil
}

func (b *bdecoder) decodeList() ([]interface{}, error) {
	var out []interface{}
	for {
		c, err := b.ReadByte()

		// if nothing to read, return all we have
		if err != nil {
			return out, nil
		}

		switch c {
		case 'i':
			if err = b.UnreadByte(); err != nil {
				return nil, err
			}
			i, err := b.decodeInt()
			if err != nil {
				break
			}
			out = append(out, i)
		case 'l':
			// if its a list, we don't unread so we can
			// continue processing its items
			l, err := b.decodeList()
			if err != nil {
				break
			}
			out = append(out, l)
		case 'd':

		case 'e':
			return out, nil
		default:
			if err = b.UnreadByte(); err != nil {
				return nil, err
			}
			s, err := b.decodeStr()
			if err != nil {
				break
			}
			out = append(out, s)
		}

		if err != nil {
			return nil, err
		}
	}
}

func (b *bdecoder) decodeDict() (map[string]interface{}, error) {

	return nil, nil
}

func Bencode(d interface{}) (string, error) {
	e := &bencoder{}
	return e.encode(d)
}

func Bdecode(s string) (interface{}, error) {
	r := bufio.NewReader(bytes.NewBufferString(s))
	d := bdecoder{*r}
	return d.decode()
}
