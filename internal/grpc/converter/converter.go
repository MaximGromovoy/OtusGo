package converter

import (
	currencyv1 "OtusGo/internal/grpc/pb/api/proto/currency/v1"
	"OtusGo/internal/model/currency"
)

// CurrencyToProto конвертирует внутренний тип валюты в protobuf сообщение
func CurrencyToProto(curr currency.CurrencyInterface) *currencyv1.Currency {
	if curr == nil {
		return nil
	}

	return &currencyv1.Currency{
		Id:    int32(curr.GetID()),
		Name:  curr.GetName(),
		Code:  curr.GetCode(),
		Value: curr.GetValue(),
	}
}

// CurrenciesToProto конвертирует массив валют в protobuf сообщения
func CurrenciesToProto(currencies []currency.CurrencyInterface) []*currencyv1.Currency {
	result := make([]*currencyv1.Currency, 0, len(currencies))
	for _, curr := range currencies {
		if curr != nil {
			result = append(result, CurrencyToProto(curr))
		}
	}
	return result
}
