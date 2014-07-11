package byteslego

import (
	"fmt"
	"sort"
	"strings"
)

type Bencoder struct{}

func (b *Bencoder) EncStr(s string) string {
	return fmt.Sprintf("%d:%s", len(s), s)
}

func (b *Bencoder) EncInt(i int64) string {
	return fmt.Sprintf("i%de", i)
}

func (b *Bencoder) EncDict(d map[string]interface{}) string {
	var encStrs []string

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
		enckey := b.EncStr(k)
		switch v.(type) {
		case string:
			encStrs = append(
				encStrs,
				fmt.Sprintf(
					"%s%s",
					enckey,
					b.EncStr(v.(string)),
				),
			)
		case int, int32, int64:
			encStrs = append(
				encStrs,
				fmt.Sprintf(
					"%s%s",
					enckey,
					b.EncInt(int64(v.(int))),
				),
			)
		case []interface{}:
			encStrs = append(
				encStrs,
				fmt.Sprintf(
					"%s%s",
					enckey,
					b.EncList(v.([]interface{})),
				),
			)
		case map[string][]interface{}:
			var sEncStrs []string
			val := v.(map[string][]interface{})
			for sk, sv := range val {
				sEncStrs = append(
					sEncStrs,
					fmt.Sprintf(
						"%s%s",
						b.EncStr(sk),
						b.EncList(sv),
					),
				)
			}

			encStrs = append(
				encStrs,
				fmt.Sprintf(
					"%sd%se",
					enckey,
					strings.Join(sEncStrs, ""),
				),
			)
		case map[string]interface{}:
			encStrs = append(
				encStrs,
				fmt.Sprintf(
					"%s%s",
					enckey,
					b.EncDict(v.(map[string]interface{})),
				),
			)
		}
	}

	return fmt.Sprintf("d%se", strings.Join(encStrs, ""))
}

func (b *Bencoder) EncList(l []interface{}) string {
	var encStrs []string
	for _, s := range l {
		switch s.(type) {
		case string:
			encStrs = append(encStrs, b.EncStr(s.(string)))
		case int, int32, int64:
			encStrs = append(encStrs, b.EncInt(int64(s.(int))))
		case []interface{}:
			encStrs = append(encStrs, b.EncList(s.([]interface{})))
		case map[string]interface{}:
			encStrs = append(encStrs, b.EncDict(s.(map[string]interface{})))
		}
	}
	return fmt.Sprintf("l%se", strings.Join(encStrs, ""))
}
