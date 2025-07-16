package models

import (
	"gophkeeper/internal/errs"
	pb "gophkeeper/internal/protos/items"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ItemPbToModels(i *pb.Item) (*Item, error) {
	item := Item{
		UserLogin: i.UserLogin,
		Name:      i.Name,
		Type:      ItemTypePbToModel(i.Type),
		Meta:      Meta{Map: i.Meta},
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

func credentialsPbToModels(cr *pb.Credentials) *Credentials {
	return &Credentials{
		Login:    cr.Login,
		Password: cr.Password,
	}
}

func textPbToModels(t *pb.Text) *Text {
	return &Text{Content: t.Content}
}

func binaryPbToModels(b *pb.Binary) *Binary {
	return &Binary{Content: b.Content}
}

func cardPbToModels(c *pb.Card) *Card {
	return &Card{
		Number:         c.Number,
		ExpiryDate:     c.ExpiryDate,
		SecurityCode:   c.SecurityCode,
		CardholderName: c.CardholderName,
	}
}

func ItemTypePbToModel(t pb.ItemType) ItemType {
	switch t {
	case pb.ItemType_ITEM_TYPE_CREDENTIALS:
		return ItemTypeCREDENTIALS
	case pb.ItemType_ITEM_TYPE_TEXT:
		return ItemTypeTEXT
	case pb.ItemType_ITEM_TYPE_BINARY:
		return ItemTypeBINARY
	case pb.ItemType_ITEM_TYPE_CARD:
		return ItemTypeCARD
	default: 
		return ItemType("unknown")
	}
}

func (i *Item) ToPb() (*pb.Item, error) {
	item := pb.Item{
		Id:        i.ID,
		UserLogin: i.UserLogin,
		Name:      i.Name,
		Type:      i.Type.ToPb(),
		Meta:      i.Meta.Map,
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
	}

	switch i.Type {
	case ItemTypeCREDENTIALS:
		item.Data = &pb.Item_Credentials{
			Credentials: i.Data.(*Credentials).ToPb(),
		}
	case ItemTypeTEXT:
		item.Data = &pb.Item_Text{
			Text: i.Data.(*Text).ToPb(),
		}
	case ItemTypeBINARY:
		item.Data = &pb.Item_Binary{
			Binary: i.Data.(*Binary).ToPb(),
		}
	case ItemTypeCARD:
		item.Data = &pb.Item_Card{
			Card: i.Data.(*Card).ToPb(),
		}
	default:
		return nil, status.Error(codes.Internal, errs.ErrIncorrectItemType.Error())
	}
	return &item, nil
}

func (cr *Credentials) ToPb() *pb.Credentials {
	return &pb.Credentials{
		Login:    cr.Login,
		Password: cr.Password,
	}
}

func (t *Text) ToPb() *pb.Text {
	return &pb.Text{Content: t.Content}
}

func (b *Binary) ToPb() *pb.Binary {
	return &pb.Binary{Content: b.Content}
}

func (c *Card) ToPb() *pb.Card {
	return &pb.Card{
		Number:         c.Number,
		ExpiryDate:     c.ExpiryDate,
		SecurityCode:   c.SecurityCode,
		CardholderName: c.CardholderName,
	}
}

func (t *ItemType) ToPb() pb.ItemType {
	switch *t {
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

