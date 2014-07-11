package byteslego

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBencoder(t *testing.T) {
	Convey("should encode string", t, func() {
		s, _ := Bencode("hello")
		So(s, ShouldEqual, "5:hello")
	})

	Convey("should encode integer", t, func() {
		s, _ := Bencode(28)
		So(s, ShouldEqual, "i28e")
	})

	Convey("should encode dictionary", t, func() {
		Convey("with string values", func() {
			s, _ := Bencode(map[string]interface{}{
				"username": "bob",
				"password": "secret",
			})
			So(s, ShouldEqual, "d8:password6:secret8:username3:bobe")
		})

		Convey("with integer values", func() {
			s, _ := Bencode(map[string]interface{}{
				"age":    30,
				"number": 38,
			})
			So(s, ShouldEqual, "d3:agei30e6:numberi38ee")
		})

		Convey("with list values", func() {
			s, _ := Bencode(map[string]interface{}{
				"fruits": []interface{}{"apple", "banana"},
				"number": []interface{}{1, 2, 3},
			})
			So(s, ShouldEqual, "d6:fruitsl5:apple6:bananae6:numberli1ei2ei3eee")
		})

		Convey("with dictionary values", func() {
			s, _ := Bencode(map[string]interface{}{
				"collection": map[string]interface{}{
					"fruit":  "apple",
					"number": 1,
				},
			})
			So(s, ShouldEqual, "d10:collectiond5:fruit5:apple6:numberi1eee")
		})

		Convey("with mixed values", func() {
			s, _ := Bencode(map[string]interface{}{
				"collection": map[string][]interface{}{
					"fruits":  []interface{}{"apple", "banana"},
					"numbers": []interface{}{1, 2, 3},
				},
			})
			So(s, ShouldEqual, "d10:collectiond6:fruitsl5:apple6:bananae7:numbersli1ei2ei3eeee")
		})
	})

	Convey("should encode list", t, func() {
		Convey("with strings", func() {
			s, _ := Bencode([]interface{}{"apple", "banana"})
			So(s, ShouldEqual, "l5:apple6:bananae")
		})

		Convey("with integers", func() {
			s, _ := Bencode([]interface{}{1, 2, 3})
			So(s, ShouldEqual, "li1ei2ei3ee")
		})

		Convey("with list", func() {
			s, _ := Bencode([]interface{}{
				[]interface{}{"chocolate", 4},
				[]interface{}{"candy", 5},
			})
			So(s, ShouldEqual, "ll9:chocolatei4eel5:candyi5eee")
		})

		Convey("with dictionary", func() {
			s, _ := Bencode([]interface{}{
				map[string]interface{}{
					"color":  "red",
					"number": 1,
				},
			})
			So(s, ShouldEqual, "ld5:color3:red6:numberi1eee")
		})

		Convey("with mixed", func() {
			s, _ := Bencode([]interface{}{
				map[string]interface{}{
					"random": map[string]interface{}{
						"colors": []interface{}{"red", "blue"},
						"number": 1,
					},
				},
			})
			So(s, ShouldEqual, "ld6:randomd6:colorsl3:red4:bluee6:numberi1eeee")
		})
	})
}

func TestBdecoder(t *testing.T) {
	Convey("should decode string", t, func() {
		Convey("valid string", func() {
			s, _ := Bdecode("5:hello")
			So(s, ShouldEqual, "hello")
		})

		Convey("string without size", func() {
			_, err := Bdecode("hello")
			So(err, ShouldNotEqual, nil)
		})

		Convey("invalid string size", func() {
			_, err := Bdecode("10:hello")
			So(err, ShouldNotEqual, nil)
		})
	})

	Convey("should decode integer", t, func() {
		Convey("valid integer", func() {
			i, _ := Bdecode("i1e")
			So(i, ShouldEqual, 1)
		})

		Convey("invalid integer", func() {
			_, err := Bdecode("20")
			So(err, ShouldNotEqual, nil)
		})
	})

	Convey("should decode list", t, func() {
		Convey("with integers", func() {
			l, _ := Bdecode("li1ei2ei3ee")
			So(l, ShouldResemble, []interface{}{int64(1), int64(2), int64(3)})
		})

		Convey("with strings", func() {
			l, _ := Bdecode("l5:apple6:bananae")
			So(l, ShouldResemble, []interface{}{"apple", "banana"})
		})

		Convey("with list", func() {
			l, _ := Bdecode("ll9:chocolatei4eel5:candyi5eee")
			So(l, ShouldResemble, []interface{}{
				[]interface{}{"chocolate", int64(4)},
				[]interface{}{"candy", int64(5)},
			})
		})
	})
}
