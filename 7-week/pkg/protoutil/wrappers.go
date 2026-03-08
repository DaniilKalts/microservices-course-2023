package protoutil

import "google.golang.org/protobuf/types/known/wrapperspb"

func StringValuePtr(value *wrapperspb.StringValue) *string {
	if value == nil {
		return nil
	}

	str := value.GetValue()
	if str == "" {
		return nil
	}

	return &str
}
