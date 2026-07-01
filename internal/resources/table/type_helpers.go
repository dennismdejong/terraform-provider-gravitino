package table

import "github.com/gravitino/terraform-provider-gravitino/internal/models"

func typeNameToDataType(name string, length, precision, scale int64) models.DataType {
	dt := models.DataType{Type: name}

	if length > 0 {
		l := length
		dt.Length = &l
	}
	if precision > 0 {
		p := precision
		dt.Precision = &p
	}
	if scale > 0 {
		s := scale
		dt.Scale = &s
	}

	switch name {
	case "integer", "long", "double", "float", "short", "byte":
		dt.Length = nil
		dt.Precision = nil
		dt.Scale = nil
	}

	return dt
}

func columnTypeToString(dt models.DataType) string {
	return dt.Type
}
