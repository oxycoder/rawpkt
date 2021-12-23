package rawpkt_test

import (
	"testing"
	"time"

	"github.com/oxycoder/rawpkt"
	"github.com/stretchr/testify/assert"
)

type MyPacketItem struct {
	ID          int ``
	Name        [5]byte
	Description string
	PosX        float32
	Uint16Test  uint16
	BoolTest    bool
	Int8Test    int8
	Int16Test   int16
	UintTest    uint
}

type MyPacketStruct struct {
	Now          time.Time
	Name         [5]byte
	Abc          int16
	OneItem      MyPacketItem
	Items        [3]MyPacketItem
	ItemSlice    []MyPacketItem
	ItemPtrSlice []*MyPacketItem
	MyString     string
}

func TestUnmarshall(t *testing.T) {
	pkt := rawpkt.NewPacket(1, true, false)
	item := MyPacketItem{
		ID:          1,
		Name:        [5]byte{0, 1, 2, 3, 4},
		Description: "Hello world",
		PosX:        3.55,
		Uint16Test:  23,
		BoolTest:    false,
		Int8Test:    3,
		Int16Test:   32,
		UintTest:    8,
	}
	items := [3]MyPacketItem{
		item, item, item,
	}
	itemptr := make([]*MyPacketItem, 0)
	pkItem := MyPacketStruct{
		Now:          time.Now(),
		Name:         [5]byte{6, 6, 6, 6, 6},
		Abc:          20,
		OneItem:      item,
		Items:        items,
		ItemSlice:    append(items[:], item),
		ItemPtrSlice: append(itemptr, &item),
		MyString:     "this is test string",
	}
	pkt.Marshal(&pkItem)
	myPkt := MyPacketStruct{}
	pkt.Unmarshal(&myPkt)
	assert.Equal(t, myPkt.Now.Unix(), pkItem.Now.Unix(), "Time should be equal, %d", myPkt.Now.Unix())
	assert.Equal(t, myPkt.Name, pkItem.Name, "Name should equal")
	assert.Equal(t, myPkt.MyString, pkItem.MyString, "MyString should equal")
	assert.Equal(t, myPkt.Abc, pkItem.Abc, "Abc should equal")
	assert.Equal(t, len(myPkt.Items), len(pkItem.Items), "Items len should equal")
	assert.Equal(t, myPkt.Items[0], pkItem.Items[0], "First item should equal")
	assert.Equal(t, myPkt.ItemPtrSlice[0].ID, pkItem.ItemPtrSlice[0].ID, "ID in ptr slice should be equal")
	assert.Equal(t, myPkt.OneItem.UintTest, pkItem.OneItem.UintTest, "Uint should be equal")
	assert.Equal(t, myPkt.Items[0].Description, pkItem.Items[0].Description, "Description should be equal")
}
