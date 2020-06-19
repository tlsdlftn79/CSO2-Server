package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

const (
	FavoriteSetLoadout   = 0
	FavoriteSetCosmetics = 1
)

type userInventory struct {
	numOfItem uint16              //物品数量
	items     []userInventoryItem //物品

	CTModel   uint32 //当前的CT模型
	TModel    uint32 //当前的T模型
	headItem  uint32 //当前的头部装饰
	gloveItem uint32 //当前的手套
	backItem  uint32 //当前的背部物品
	stepsItem uint32 //当前的脚步效果
	cardItem  uint32 //当前的卡片
	sprayItem uint32 //当前的喷漆

	buyMenu  userBuyMenu //购买菜单
	loadouts []userLoadout
}
type userInventoryItem struct {
	id    uint32 //物品id
	count uint16 //数量
}

type inFavoritePacket struct {
	packetType uint8
}

func onFavorite(seq *uint8, p packet, client net.Conn) {
	var pkt inFavoritePacket
	if !praseFavoritePacket(p, &pkt) {
		log.Println("Error : Recived a illegal favorite packet from", client.RemoteAddr().String())
		return
	}
	switch pkt.packetType {
	case FavoriteSetLoadout:
		//log.Println("Recived a favorite SetLoadout packet from", client.RemoteAddr().String())
		onFavoriteSetLoadout(seq, p, client)
	case FavoriteSetCosmetics:
		//log.Println("Recived a favorite SetCosmetics packet from", client.RemoteAddr().String())
		onFavoriteSetCosmetics(seq, p, client)
	default:
		log.Println("Unknown favorite packet", pkt.packetType, "from", client.RemoteAddr().String())
	}
}

func praseFavoritePacket(p packet, dest *inFavoritePacket) bool {
	if p.datalen < 6 {
		return false
	}
	offset := 5
	(*dest).packetType = ReadUint8(p.data, &offset)
	return true
}

func BuildInventoryInfo(u user) []byte {
	buf := make([]byte, 5+u.inventory.numOfItem*11)
	offset := 0
	WriteUint16(&buf, u.inventory.numOfItem, &offset)
	for k, v := range u.inventory.items {
		WriteUint16(&buf, uint16(k), &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint32(&buf, v.id, &offset)
		WriteUint16(&buf, v.count, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 0, &offset)
	}
	return buf
}

func BuildUnlockReply() []byte {
	buf := []byte{0x01, 0x4B, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x0B, 0x00, 0x00, 0x00, 0x01, 0xE8, 0x03, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x0C, 0x00,
		0x00, 0x00, 0x01, 0xDC, 0x05, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x01,
		0xE8, 0x03, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x01, 0xDC, 0x05, 0x00,
		0x00, 0x0B, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x00, 0x00, 0x01, 0x08, 0x07, 0x00, 0x00, 0x3C, 0x00,
		0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x01, 0x80, 0xBB, 0x00, 0x00, 0x1F, 0x00, 0x00, 0x00, 0x11,
		0x00, 0x00, 0x00, 0x01, 0xC0, 0x5D, 0x00, 0x00, 0x11, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00,
		0x01, 0x08, 0x07, 0x00, 0x00, 0x1C, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x00, 0x01, 0x4C, 0x1D,
		0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x01, 0x60, 0x61, 0x02, 0x00, 0x35,
		0x00, 0x00, 0x00, 0x15, 0x00, 0x00, 0x00, 0x01, 0x30, 0x75, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00,
		0x16, 0x00, 0x00, 0x00, 0x01, 0xA0, 0x0F, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00, 0x17, 0x00, 0x00,
		0x00, 0x01, 0x98, 0x3A, 0x00, 0x00, 0x3F, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x01, 0xE0,
		0x93, 0x04, 0x00, 0x14, 0x00, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00, 0x01, 0xA0, 0x0F, 0x00, 0x00,
		0x07, 0x00, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x01, 0x98, 0x3A, 0x00, 0x00, 0x3E, 0x00, 0x00,
		0x00, 0x1B, 0x00, 0x00, 0x00, 0x01, 0xE0, 0x93, 0x04, 0x00, 0x05, 0x00, 0x00, 0x00, 0x1C, 0x00,
		0x00, 0x00, 0x01, 0x08, 0x07, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x1D, 0x00, 0x00, 0x00, 0x01,
		0x30, 0x75, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x1E, 0x00, 0x00, 0x00, 0x01, 0x88, 0x13, 0x00,
		0x00, 0x0C, 0x00, 0x00, 0x00, 0x1F, 0x00, 0x00, 0x00, 0x01, 0x20, 0x4E, 0x00, 0x00, 0x16, 0x00,
		0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x01, 0x20, 0x4E, 0x00, 0x00, 0x34, 0x00, 0x00, 0x00, 0x43,
		0x00, 0x00, 0x00, 0x01, 0x30, 0x75, 0x00, 0x00, 0x46, 0x00, 0x00, 0x00, 0x57, 0x00, 0x00, 0x00,
		0x01, 0x20, 0xA1, 0x07, 0x00, 0x47, 0x00, 0x00, 0x00, 0x58, 0x00, 0x00, 0x00, 0x01, 0x20, 0xA1,
		0x07, 0x00, 0x4D, 0x00, 0x00, 0x00, 0x59, 0x00, 0x00, 0x00, 0x00, 0x90, 0x01, 0x00, 0x00, 0x55,
		0x00, 0x00, 0x00, 0x81, 0x00, 0x00, 0x00, 0x00, 0x70, 0x03, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00,
		0x90, 0x00, 0x00, 0x00, 0x01, 0x30, 0x75, 0x00, 0x00, 0x1D, 0x00, 0x00, 0x00, 0x91, 0x00, 0x00,
		0x00, 0x01, 0x60, 0xEA, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x92, 0x00, 0x00, 0x00, 0x01, 0x48,
		0xE8, 0x01, 0x00, 0x2F, 0x00, 0x00, 0x00, 0x93, 0x00, 0x00, 0x00, 0x01, 0x40, 0x0D, 0x03, 0x00,
		0x6A, 0xBF, 0x00, 0x00, 0xA8, 0x00, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x70, 0xBF, 0x00,
		0x00, 0xA9, 0x00, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00, 0x6F, 0xBF, 0x00, 0x00, 0xAA, 0x00,
		0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x6E, 0xBF, 0x00, 0x00, 0xAB, 0x00, 0x00, 0x00, 0x00,
		0x50, 0x00, 0x00, 0x00, 0x69, 0xBF, 0x00, 0x00, 0xAC, 0x00, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00,
		0x00, 0x72, 0xBF, 0x00, 0x00, 0xAD, 0x00, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00, 0x6B, 0xBF,
		0x00, 0x00, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x6D, 0xBF, 0x00, 0x00, 0xAF,
		0x00, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00, 0x4A, 0x00, 0x00, 0x00, 0xD7, 0x00, 0x00, 0x00,
		0x01, 0x50, 0xC3, 0x00, 0x00, 0x4B, 0x00, 0x00, 0x00, 0xD8, 0x00, 0x00, 0x00, 0x01, 0x00, 0x77,
		0x01, 0x00, 0x4E, 0x00, 0x00, 0x00, 0xE8, 0x00, 0x00, 0x00, 0x01, 0x70, 0x11, 0x01, 0x00, 0x52,
		0x00, 0x00, 0x00, 0xE9, 0x00, 0x00, 0x00, 0x01, 0xC0, 0xD4, 0x01, 0x00, 0x5B, 0x00, 0x00, 0x00,
		0x06, 0x01, 0x00, 0x00, 0x01, 0xF0, 0x49, 0x02, 0x00, 0x5F, 0x00, 0x00, 0x00, 0x19, 0x01, 0x00,
		0x00, 0x01, 0x60, 0xEA, 0x00, 0x00, 0x60, 0x00, 0x00, 0x00, 0x1A, 0x01, 0x00, 0x00, 0x01, 0xC0,
		0xD4, 0x01, 0x00, 0x64, 0x00, 0x00, 0x00, 0x38, 0x01, 0x00, 0x00, 0x01, 0xF0, 0x49, 0x02, 0x00,
		0x68, 0x00, 0x00, 0x00, 0x5C, 0x01, 0x00, 0x00, 0x01, 0x20, 0xA1, 0x07, 0x00, 0x6D, 0x00, 0x00,
		0x00, 0x82, 0x01, 0x00, 0x00, 0x01, 0xA0, 0x86, 0x01, 0x00, 0x6C, 0x00, 0x00, 0x00, 0x83, 0x01,
		0x00, 0x00, 0x01, 0xA0, 0x86, 0x01, 0x00, 0x6E, 0x00, 0x00, 0x00, 0x84, 0x01, 0x00, 0x00, 0x01,
		0xA0, 0x86, 0x01, 0x00, 0x42, 0x00, 0x00, 0x00, 0xFA, 0x01, 0x00, 0x00, 0x01, 0x30, 0x75, 0x00,
		0x00, 0x43, 0x00, 0x00, 0x00, 0xFB, 0x01, 0x00, 0x00, 0x01, 0x50, 0xC3, 0x00, 0x00, 0x78, 0x00,
		0x00, 0x00, 0xFC, 0x01, 0x00, 0x00, 0x01, 0x40, 0x0D, 0x03, 0x00, 0x79, 0x00, 0x00, 0x00, 0x07,
		0x02, 0x00, 0x00, 0x00, 0xA0, 0x00, 0x00, 0x00, 0x7C, 0x00, 0x00, 0x00, 0x08, 0x02, 0x00, 0x00,
		0x00, 0x04, 0x01, 0x00, 0x00, 0x7A, 0x00, 0x00, 0x00, 0x09, 0x02, 0x00, 0x00, 0x00, 0xE0, 0x01,
		0x00, 0x00, 0x7B, 0x00, 0x00, 0x00, 0x0A, 0x02, 0x00, 0x00, 0x00, 0x44, 0x02, 0x00, 0x00, 0x7D,
		0x00, 0x00, 0x00, 0x58, 0x02, 0x00, 0x00, 0x00, 0x44, 0x02, 0x00, 0x00, 0x7E, 0x00, 0x00, 0x00,
		0x59, 0x02, 0x00, 0x00, 0x00, 0x0C, 0x03, 0x00, 0x00, 0x81, 0x00, 0x00, 0x00, 0x91, 0x02, 0x00,
		0x00, 0x01, 0xF0, 0x49, 0x02, 0x00, 0x82, 0x00, 0x00, 0x00, 0x92, 0x02, 0x00, 0x00, 0x01, 0x00,
		0x53, 0x07, 0x00, 0x83, 0x00, 0x00, 0x00, 0x93, 0x02, 0x00, 0x00, 0x01, 0x60, 0x5B, 0x03, 0x00,
		0x85, 0x00, 0x00, 0x00, 0x94, 0x02, 0x00, 0x00, 0x00, 0x40, 0x01, 0x00, 0x00, 0x84, 0x00, 0x00,
		0x00, 0x95, 0x02, 0x00, 0x00, 0x00, 0x08, 0x02, 0x00, 0x00, 0x87, 0x00, 0x00, 0x00, 0x1F, 0x03,
		0x00, 0x00, 0x00, 0x08, 0x02, 0x00, 0x00, 0x8A, 0x00, 0x00, 0x00, 0xA4, 0x03, 0x00, 0x00, 0x01,
		0xE0, 0x93, 0x04, 0x00, 0x8F, 0x00, 0x00, 0x00, 0x44, 0x04, 0x00, 0x00, 0x01, 0x80, 0xA9, 0x03,
		0x00, 0x90, 0x00, 0x00, 0x00, 0x45, 0x04, 0x00, 0x00, 0x01, 0x40, 0x7E, 0x05, 0x00, 0x91, 0x00,
		0x00, 0x00, 0x46, 0x04, 0x00, 0x00, 0x01, 0x00, 0x53, 0x07, 0x00, 0x9B, 0x00, 0x00, 0x00, 0xA9,
		0x04, 0x00, 0x00, 0x01, 0xF0, 0x49, 0x02, 0x00, 0x9C, 0x00, 0x00, 0x00, 0xAA, 0x04, 0x00, 0x00,
		0x01, 0x40, 0x0D, 0x03, 0x00, 0x97, 0x00, 0x00, 0x00, 0xFC, 0x04, 0x00, 0x00, 0x01, 0x42, 0x99,
		0x00, 0x00, 0x98, 0x00, 0x00, 0x00, 0xFD, 0x04, 0x00, 0x00, 0x01, 0x86, 0x29, 0x02, 0x00, 0x99,
		0x00, 0x00, 0x00, 0xFE, 0x04, 0x00, 0x00, 0x01, 0x8C, 0xED, 0x02, 0x00, 0x10, 0x00, 0x03, 0x00,
		0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x42, 0x00, 0x00, 0x00, 0x43, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x07, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x00, 0x00,
		0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x00, 0x00, 0x00, 0x1C, 0x00,
		0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x00, 0x00, 0x00,
		0x35, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x00,
		0x00, 0x00, 0x34, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x13, 0x00, 0x00, 0x00, 0x4D, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x13, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x3E, 0x00, 0x00, 0x00, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0x00, 0x00, 0x00, 0x11, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x3F, 0x00,
		0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00,
		0x19, 0x00, 0x00, 0x00, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00,
		0x06, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00,
		0x0B, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x00, 0x00,
		0x10, 0x00, 0x00, 0x00, 0x11, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x00,
		0x14, 0x00, 0x00, 0x00, 0x15, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00,
		0x1A, 0x00, 0x00, 0x00, 0x1C, 0x00, 0x00, 0x00, 0x6C, 0xBF, 0x00, 0x00, 0x71, 0xBF, 0x00, 0x00,
		0x42, 0x00, 0x00, 0x00, 0x94, 0x01, 0x00, 0x00}
	return buf
}

func createNewUserInventory() userInventory {
	return userInventory{
		25,
		createDeafaultInventoryItem(),
		1047,
		1048,
		0,
		0,
		0,
		0,
		0,
		42001,
		createNewUserBuyMenu(),
		createNewLoadout(),
	}
}

func createDeafaultInventoryItem() []userInventoryItem {
	items := []userInventoryItem{}
	number := []uint32{2, 3, 4, 6, 8, 13, 14, 15, 18, 19, 21, 23, 27, 34, 36, 37, 80, 128, 101, 1001, 1002, 1003, 1004, 49009, 49004}
	for _, v := range number {
		items = append(items, userInventoryItem{v, 1})
	}
	return items
}