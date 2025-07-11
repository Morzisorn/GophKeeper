package controllers

import (
	"context"
	"gophkeeper/internal/errs"
	pb "gophkeeper/internal/protos/items"
	iserv "gophkeeper/internal/server/services/item_service"
	"gophkeeper/models"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ItemController struct {
	pb.UnimplementedItemsServiceServer
	service *iserv.ItemService
}

func NewItemController(service *iserv.ItemService) *ItemController {
	return &ItemController{
		service: service,
	}
}

func (ic *ItemController) AddItem(ctx context.Context, in *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	if !isItemValid(in.Item) {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	item, err := itemPbToModels(in.Item)
	if err != nil {
		switch err {
		case errs.ErrIncorrectItemType:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, errs.ErrInternalServerError.Error())
		}
	}

	err = ic.service.AddItem(ctx, item)
	if err != nil {
		switch err {
		case errs.ErrItemAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, errs.ErrInternalServerError.Error())
		}
	}
	return &pb.AddItemResponse{
		Success: true,
	}, nil
}

func (ic *ItemController) EditItem(ctx context.Context, in *pb.EditItemRequest) (*pb.EditItemResponse, error) {
	if !isItemValid(in.Item) {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	item, err := itemPbToModels(in.Item)
	if err != nil {
		switch err {
		case errs.ErrIncorrectItemType:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, errs.ErrInternalServerError.Error())
		}
	}

	err = ic.service.EditItem(ctx, item)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EditItemResponse{
		Success: true,
	}, nil

}

func isItemValid(i *pb.Item) bool {
	return i.Name != "" && i.Type.String() != "" && i.UserLogin != "" && i.Data != nil
}

func itemPbToModels(i *pb.Item) (*models.Item, error) {
	item := models.Item{
		UserLogin: i.UserLogin,
		Name:      i.Name,
		Type:      models.ItemType(i.Type.String()),
		Meta:      models.Meta{Map: i.Meta},
	}

	switch i.Type {
	case pb.ItemType_ITEM_TYPE_CREDENTIALS:
		item.Data = credentialsPbToModels(i.GetCredentials())
	case pb.ItemType_ITEM_TYPE_TEXT:
		item.Data = textPbToModels(i.GetText())
	case pb.ItemType_ITEM_TYPE_BINARY:
		item.Data = binaryPbToModels(i.GetBinary())
	case pb.ItemType_ITEM_TYPE_CARD:
		item.Data = cardPbToModels(i.GetCard())
	default:
		return nil, errs.ErrIncorrectItemType
	}
	return &item, nil
}

func credentialsPbToModels(cr *pb.Credentials) *models.Credentials {
	return &models.Credentials{
		Login:    cr.Login,
		Password: cr.Password,
	}
}

func textPbToModels(t *pb.Text) *models.Text {
	return &models.Text{Content: t.Content}
}

func binaryPbToModels(b *pb.Binary) *models.Binary {
	return &models.Binary{Content: b.Content}
}

func cardPbToModels(c *pb.Card) *models.Card {
	return &models.Card{
		Number:         c.Number,
		ExpiryDate:     c.ExpiryDate,
		SecurityCode:   c.SecurityCode,
		CardholderName: c.CardholderName,
	}
}

func (ic *ItemController) DeleteItem(ctx context.Context, in *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	if in.ItemId == "" || in.UserLogin == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	err := ic.service.DeleteItem(ctx, in.UserLogin, in.ItemId)
	if err != nil {
		switch err {
		case errs.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case errs.ErrItemNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.DeleteItemResponse{
		Success: true,
	}, nil
}

func (ic *ItemController) GetItems(ctx context.Context, in *pb.GetUserItemsRequest) (*pb.GetUserItemsResponse, error) {
	if in.UserLogin == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	items, err := ic.service.GetUserItems(ctx, models.ItemType(in.Type.String()), in.UserLogin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbItems := make(map[string]*pb.Item, len(items))
	for i, item := range items {
		pbItem, err := itemModelsToPb(&item)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		pbItems[i] = pbItem
	}

	return &pb.GetUserItemsResponse{
		Items: pbItems,
	}, nil
}

func itemModelsToPb(i *models.Item) (*pb.Item, error) {
	item := pb.Item{
		Id:        i.ID,
		UserLogin: i.UserLogin,
		Name:      i.Name,
		Type:      pb.ItemType(pb.ItemType_value[strings.ToUpper((i.Type.String()))]),
		Meta:      i.Meta.Map,
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
	}

	switch i.Type {
	case models.ItemTypeCREDENTIALS:
		item.Data = &pb.Item_Credentials{
			Credentials: credentialsModelsToPb(i.Data.(*models.Credentials)),
		}
	case models.ItemTypeTEXT:
		item.Data = &pb.Item_Text{
			Text: textModelsToPb(i.Data.(*models.Text)),
		}
	case models.ItemTypeBINARY:
		item.Data = &pb.Item_Binary{
			Binary: binaryModelsToPb(i.Data.(*models.Binary)),
		}
	case models.ItemTypeCARD:
		item.Data = &pb.Item_Card{
			Card: cardModelsToPb(i.Data.(*models.Card)),
		}
	default:
		return nil, status.Error(codes.Internal, errs.ErrIncorrectItemType.Error())
	}
	return &item, nil
}

func credentialsModelsToPb(cr *models.Credentials) *pb.Credentials {
	return &pb.Credentials{
		Login:    cr.Login,
		Password: cr.Password,
	}
}

func textModelsToPb(t *models.Text) *pb.Text {
	return &pb.Text{Content: t.Content}
}

func binaryModelsToPb(b *models.Binary) *pb.Binary {
	return &pb.Binary{Content: b.Content}
}

func cardModelsToPb(c *models.Card) *pb.Card {
	return &pb.Card{
		Number:         c.Number,
		ExpiryDate:     c.ExpiryDate,
		SecurityCode:   c.SecurityCode,
		CardholderName: c.CardholderName,
	}
}

func (ic *ItemController) GetItemTypesCounters(ctx context.Context, in *pb.TypesCountsRequest) (*pb.TypesCountsResponse, error) {
	counters, err := ic.service.GetTypesCounts(ctx, in.UserLogin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbCounters := make(map[string]int32, len(counters))

	for k, v := range counters {
		pbCounters[k.String()] = v
	}

	return &pb.TypesCountsResponse{
		Types: pbCounters,
	}, nil
}