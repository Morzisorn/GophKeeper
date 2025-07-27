package controllers

import (
	"context"
	"gophkeeper/internal/errs"
	pb "gophkeeper/internal/protos/items"
	iserv "gophkeeper/internal/server/services/item_service"
	"gophkeeper/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ItemController struct {
	pb.UnimplementedItemsControllerServer
	service *iserv.ItemService
}

func NewItemController(service *iserv.ItemService) *ItemController {
	return &ItemController{
		service: service,
	}
}

func (ic *ItemController) AddItem(ctx context.Context, in *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "login not found in context")
	}
	if login != in.Item.UserLogin {
		return nil, status.Error(codes.Unauthenticated, "login in JWT differ from item user login")
	}

	if !isPbItemValid(in.Item) {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	item := models.EncryptedItemPbToModels(in.Item)

	if err := ic.service.AddItem(ctx, item); err != nil {
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
	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "login not found in context")
	}
	if login != in.Item.UserLogin {
		return nil, status.Error(codes.Unauthenticated, "login in JWT differ from item user login")
	}

	if !isPbItemValid(in.Item) {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	item := models.EncryptedItemPbToModels(in.Item)

	if err := ic.service.EditItem(ctx, item); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EditItemResponse{
		Success: true,
	}, nil

}

func isPbItemValid(i *pb.EncryptedItem) bool {
	return i.Name != "" && i.Type.String() != "" && i.UserLogin != "" && i.EncryptedData.EncryptedContent != "" && i.EncryptedData.Nonce != ""
}

func (ic *ItemController) DeleteItem(ctx context.Context, in *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "login not found in context")
	}
	if login != in.UserLogin {
		return nil, status.Error(codes.Unauthenticated, "login in JWT differ from item user login")
	}

	if in.ItemId == nil || in.UserLogin == "" {
		return nil, status.Error(codes.InvalidArgument, errs.ErrRequiredArgumentIsMissing.Error())
	}

	err := ic.service.DeleteItem(ctx, in.UserLogin, models.ItemIdPbToModels(in.ItemId))
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

func (ic *ItemController) GetUserItems(ctx context.Context, in *pb.GetUserItemsRequest) (*pb.GetUserItemsResponse, error) {
	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "login not found in context")
	}
	if login != in.UserLogin {
		return nil, status.Error(codes.Unauthenticated, "login in JWT differ from item user login")
	}

	items, err := ic.service.GetUserItems(ctx, models.ItemTypePbToModel(in.Type), in.UserLogin)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbItems := make([]*pb.EncryptedItem, len(items))
	for i, item := range items {
		pbItem, err := item.ToPb()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		pbItems[i] = pbItem
	}

	return &pb.GetUserItemsResponse{
		Items: pbItems,
	}, nil
}

func (ic *ItemController) GetItemTypesCounters(ctx context.Context, in *pb.TypesCountsRequest) (*pb.TypesCountsResponse, error) {
	login, ok := ctx.Value("login").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "login not found in context")
	}
	if login != in.UserLogin {
		return nil, status.Error(codes.Unauthenticated, "login in JWT differ from item user login")
	}
	
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
