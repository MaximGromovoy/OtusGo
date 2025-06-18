package server

import (
	"OtusGo/internal/grpc/converter"
	currencyv1 "OtusGo/internal/grpc/pb/api/proto/currency/v1"
	"OtusGo/internal/model/currency"
	"OtusGo/internal/repository/currencyRepository"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyGRPCServer struct {
	currencyv1.UnimplementedCurrencyServiceServer
	repo *currencyRepository.CurrencyRepository
}

func NewCurrencyGRPCServer(repo *currencyRepository.CurrencyRepository) *CurrencyGRPCServer {
	return &CurrencyGRPCServer{
		repo: repo,
	}
}

func (s *CurrencyGRPCServer) CreateCurrency(ctx context.Context, req *currencyv1.CreateCurrencyRequest) (*currencyv1.CreateCurrencyResponse, error) {
	if !currency.IsCurrencySupported(req.CurrencyType) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency type: %s", req.CurrencyType)
	}

	newCurrency := currency.NewCurrency(req.CurrencyType, req.Value)
	if newCurrency == nil {
		return nil, status.Errorf(codes.Internal, "failed to create currency")
	}

	err := s.repo.Add(newCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add currency: %v", err)
	}

	return &currencyv1.CreateCurrencyResponse{
		Currency: converter.CurrencyToProto(newCurrency),
		Message:  fmt.Sprintf("Currency %s created successfully", req.CurrencyType),
	}, nil
}

func (s *CurrencyGRPCServer) GetCurrency(ctx context.Context, req *currencyv1.GetCurrencyRequest) (*currencyv1.GetCurrencyResponse, error) {
	if !currency.IsCurrencySupported(req.CurrencyType) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency type: %s", req.CurrencyType)
	}

	curr, err := s.repo.GetByTypeAndID(req.CurrencyType, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "currency not found: %v", err)
	}

	return &currencyv1.GetCurrencyResponse{
		Currency: converter.CurrencyToProto(curr),
	}, nil
}

func (s *CurrencyGRPCServer) GetCurrencies(ctx context.Context, req *currencyv1.GetCurrenciesRequest) (*currencyv1.GetCurrenciesResponse, error) {
	if !currency.IsCurrencySupported(req.CurrencyType) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency type: %s", req.CurrencyType)
	}

	currencies := s.repo.GetAll(req.CurrencyType)

	return &currencyv1.GetCurrenciesResponse{
		Currencies: converter.CurrenciesToProto(currencies),
	}, nil
}

func (s *CurrencyGRPCServer) UpdateCurrency(ctx context.Context, req *currencyv1.UpdateCurrencyRequest) (*currencyv1.UpdateCurrencyResponse, error) {
	if !currency.IsCurrencySupported(req.CurrencyType) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency type: %s", req.CurrencyType)
	}

	updatedCurrency := currency.NewCurrencyWithID(req.CurrencyType, req.Value, int(req.Id))
	if updatedCurrency == nil {
		return nil, status.Errorf(codes.Internal, "failed to create currency for update")
	}

	err := s.repo.Update(req.CurrencyType, int(req.Id), updatedCurrency)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to update currency: %v", err)
	}

	return &currencyv1.UpdateCurrencyResponse{
		Currency: converter.CurrencyToProto(updatedCurrency),
		Message:  fmt.Sprintf("Currency %s with ID %d updated successfully", req.CurrencyType, req.Id),
	}, nil
}

func (s *CurrencyGRPCServer) DeleteCurrency(ctx context.Context, req *currencyv1.DeleteCurrencyRequest) (*currencyv1.DeleteCurrencyResponse, error) {
	if !currency.IsCurrencySupported(req.CurrencyType) {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency type: %s", req.CurrencyType)
	}

	err := s.repo.DeleteByTypeAndID(req.CurrencyType, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete currency: %v", err)
	}

	return &currencyv1.DeleteCurrencyResponse{
		Message: fmt.Sprintf("Currency %s with ID %d deleted successfully", req.CurrencyType, req.Id),
	}, nil
}
