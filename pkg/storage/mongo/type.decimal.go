package mongo

// "github.com/ericlagergren/decimal"
// "go.mongodb.org/mongo-driver/bson/bsontype"
// "go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

// type mongoBig decimal.Big

// type blah bsoncore.Value // todo delete

// func (mb mongoBig) MarshalBSONValue() (bsontype.Type, []byte, error) {
// 	// return bsontype.Binary, bsoncore.AppendBinary(nil, 4, mb[:]), nil
// 	return bsontype.Decimal128, bsoncore.AppendDecimal128(), nil
// }

// func (mb *mongoBig) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
// 	if t != bsontype.Decimal128 {
// 		return fmt.Errorf("invalid format on unmarshal bson value")
// 	}

// data, _, ok := bsoncore.ReadDecimal128(raw)
// if !ok {
// 	return fmt.Errorf("not enough bytes to unmarshal bson value")
// }

// copy(mb, data)

// 	return nil
// }
