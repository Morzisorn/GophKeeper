package models

import (
	pb "gophkeeper/internal/protos/items"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func EncryptedItemPbToModels(i *pb.EncryptedItem) *EncryptedItem {
	return &EncryptedItem{
		ID:            ItemIdPbToModels(i.Id),
		UserLogin:     i.UserLogin,
		Name:          i.Name,
		Type:          ItemTypePbToModel(i.Type),
		EncryptedData: EncryptedDataPbToModel(i.EncryptedData),
		Meta:          Meta{Map: i.Meta},
		CreatedAt:     i.CreatedAt.AsTime(),
		UpdatedAt:     i.UpdatedAt.AsTime(),
	}
}

func ItemIdPbToModels(idPb []byte) [16]byte {
	var id [16]byte
	copy(id[:], idPb)
	return id
}

func (i *EncryptedItem) ToPb() (*pb.EncryptedItem, error) {
	item := pb.EncryptedItem{
		Id:            i.ID[:],
		UserLogin:     i.UserLogin,
		Name:          i.Name,
		Type:          i.Type.ToPb(),
		EncryptedData: i.EncryptedData.ToPb(),
		Meta:          i.Meta.Map,
		CreatedAt:     timestamppb.New(i.CreatedAt),
		UpdatedAt:     timestamppb.New(i.UpdatedAt),
	}

	return &item, nil
}

func (ed *EncryptedData) ToPb() *pb.EncryptedData {
	return &pb.EncryptedData{
		EncryptedContent: ed.EncryptedContent,
		Nonce:            ed.Nonce,
	}
}

func EncryptedDataPbToModel(ed *pb.EncryptedData) EncryptedData {
	return EncryptedData{
		EncryptedContent: ed.EncryptedContent,
		Nonce:            ed.Nonce,
	}
}

func ItemTypePbToModel(t pb.ItemType) ItemType {
	switch t {
	case pb.ItemType_ITEM_TYPE_UNSPECIFIED:
		return ItemTypeUNSPECIFIED
	case pb.ItemType_ITEM_TYPE_CREDENTIALS:
		return ItemTypeCREDENTIALS
	case pb.ItemType_ITEM_TYPE_TEXT:
		return ItemTypeTEXT
	case pb.ItemType_ITEM_TYPE_BINARY:
		return ItemTypeBINARY
	case pb.ItemType_ITEM_TYPE_CARD:
		return ItemTypeCARD
	default:
		return ItemTypeUNSPECIFIED
	}
}

func (t *ItemType) ToPb() pb.ItemType {
	switch *t {
	case ItemTypeUNSPECIFIED:
		return pb.ItemType_ITEM_TYPE_UNSPECIFIED
	case ItemTypeCREDENTIALS:
		return pb.ItemType_ITEM_TYPE_CREDENTIALS
	case ItemTypeTEXT:
		return pb.ItemType_ITEM_TYPE_TEXT
	case ItemTypeBINARY:
		return pb.ItemType_ITEM_TYPE_BINARY
	case ItemTypeCARD:
		return pb.ItemType_ITEM_TYPE_CARD
	default:
		return pb.ItemType_ITEM_TYPE_UNSPECIFIED
	}
}
