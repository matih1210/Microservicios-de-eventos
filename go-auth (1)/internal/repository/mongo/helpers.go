package mongo

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func oidToHex(id interface{}) string {
    if oid, ok := id.(primitive.ObjectID); ok {
        return oid.Hex()
    }
    return ""
}

func hexToOID(id string) primitive.ObjectID {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return primitive.NilObjectID
    }
    return oid
}
