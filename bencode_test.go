package byteslego

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBencoder(t *testing.T) {
	bencoder := &Bencoder{}

	Convey("should encode string", t, func() {
		So(bencoder.EncStr("hello"), ShouldEqual, "5:hello")
	})

	Convey("should encode integer", t, func() {
		So(bencoder.EncInt(28), ShouldEqual, "i28e")
	})

	Convey("should encode dictionary", t, func() {
		Convey("with string values", func() {
			So(
				bencoder.EncDict(map[string]interface{}{
					"username": "bob",
					"password": "secret",
				}),
				ShouldEqual,
				"d8:password6:secret8:username3:bobe",
			)
		})

		Convey("with integer values", func() {
			So(
				bencoder.EncDict(map[string]interface{}{
					"age":    30,
					"number": 38,
				}),
				ShouldEqual,
				"d3:agei30e6:numberi38ee",
			)
		})

		Convey("with list values", func() {
			So(
				bencoder.EncDict(map[string]interface{}{
					"fruits": []interface{}{"apple", "banana"},
					"number": []interface{}{1, 2, 3},
				}),
				ShouldEqual,
				"d6:fruitsl5:apple6:bananae6:numberli1ei2ei3eee",
			)
		})

		Convey("with dictionary values", func() {
			So(
				bencoder.EncDict(map[string]interface{}{
					"collection": map[string]interface{}{
						"fruit":  "apple",
						"number": 1,
					},
				}),
				ShouldEqual,
				"d10:collectiond5:fruit5:apple6:numberi1eee",
			)
		})

		Convey("with mixed values", func() {
			So(
				bencoder.EncDict(map[string]interface{}{
					"collection": map[string][]interface{}{
						"fruits":  []interface{}{"apple", "banana"},
						"numbers": []interface{}{1, 2, 3},
					},
				}),
				ShouldEqual,
				"d10:collectiond6:fruitsl5:apple6:bananae7:numbersli1ei2ei3eeee",
			)
		})
	})

	Convey("should encode list", t, func() {
		Convey("with strings", func() {
			So(
				bencoder.EncList([]interface{}{"apple", "banana"}),
				ShouldEqual,
				"l5:apple6:bananae",
			)
		})

		Convey("with integers", func() {
			So(
				bencoder.EncList([]interface{}{1, 2, 3}),
				ShouldEqual,
				"li1ei2ei3ee",
			)
		})

		Convey("with list", func() {
			So(
				bencoder.EncList([]interface{}{
					[]interface{}{"chocolate", 4},
					[]interface{}{"candy", 5},
				}),
				ShouldEqual,
				"ll9:chocolatei4eel5:candyi5eee",
			)
		})

		Convey("with dictionary", func() {
			So(
				bencoder.EncList([]interface{}{
					map[string]interface{}{
						"color":  "red",
						"number": 1,
					},
				}),
				ShouldEqual,
				"ld5:color3:red6:numberi1eee",
			)
		})

		Convey("with mixed", func() {
			So(
				bencoder.EncList([]interface{}{
					map[string]interface{}{
						"random": map[string]interface{}{
							"colors": []interface{}{"red", "blue"},
							"number": 1,
						},
					},
				}),
				ShouldEqual,
				"ld6:randomd6:colorsl3:red4:bluee6:numberi1eeee",
			)
		})
	})
}
