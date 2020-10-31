package main

import (
	"fmt"
	"log"

	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
)

func main() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&mainForm)
	vcl.Application.Run()
}

type TMainForm struct {
	*vcl.TForm
}

var (
	mainForm *TMainForm
)

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("PlieDriver AMD")
	f.SetWidth(800)
	f.SetHeight(600)
	f.SetOnClose(func(Sender vcl.IObject, Action *types.TCloseAction) {
		fmt.Println("close")
	})

	// 约束窗口
	f.SetOnConstrainedResize(func(sender vcl.IObject, minWidth, minHeight, maxWidth, maxHeight *int32) {
		*minWidth = 800
		*minHeight = 600
		*maxWidth = 800
		*maxHeight = 600
	})

	f.SetOnCloseQuery(func(Sender vcl.IObject, CanClose *bool) {
		*CanClose = vcl.MessageDlg("是否退出?", types.MtInformation, types.MbYes, types.MbNo) == types.MrYes
		fmt.Println("OnCloseQuery")
	})

	f.menuInit()
	f.bgLogoInit()
	f.fromInit()
}

func (f *TMainForm) bgLogoInit() {
	img := vcl.NewImage(f)
	img.SetBounds(0, 0, 800, 800)
	img.SetParent(f)
	img.Picture().LoadFromFile("./img/amd.png")
	//img.SetStretch(true)
	img.SetProportional(true)
}

func (f *TMainForm) menuInit() {
	mainMenu := vcl.NewMainMenu(f)
	item := vcl.NewMenuItem(f)
	item.SetCaption("About(&F)")
	mainMenu.Items().Add(item)

	item2 := vcl.NewMenuItem(f)
	item2.SetCaption("About")
	item2.SetOnClick(func(vcl.IObject) {
		vcl.ShowMessage("By: WorldLink 2021 AMD定制版")
	})
	item.Add(item2)

	item2 = vcl.NewMenuItem(f)
	item2.SetCaption("Exit(&E)")
	item2.SetShortCutFromString("Ctrl+Q")
	item2.SetOnClick(func(vcl.IObject) {
		mainForm.Close()
	})
	item.Add(item2)
}

func (f *TMainForm) fromInit() {
	addr := vcl.NewLabel(f)
	addr.SetCaption("地址(address): ")
	addr.SetBounds(100, 150, 32, 32)
	addr.SetParent(f)

	id := vcl.NewLabel(f)
	id.SetCaption("id(user id): ")
	id.SetBounds(100, 200, 32, 32)
	id.SetParent(f)

	passwd := vcl.NewLabel(f)
	passwd.SetCaption("密码(password): ")
	passwd.SetBounds(100, 250, 32, 32)
	passwd.SetParent(f)

	edit := vcl.NewEdit(mainForm)
	edit.SetParent(mainForm)
	edit.SetLeft(220)
	edit.SetTop(150)
	edit.SetWidth(200)
	edit.SetTextHint("address")

	edit2 := vcl.NewEdit(mainForm)
	edit2.SetParent(mainForm)
	edit2.SetLeft(220)
	edit2.SetTop(200)
	edit2.SetWidth(200)
	edit2.SetTextHint("user id")

	edit3 := vcl.NewEdit(mainForm)
	edit3.SetParent(mainForm)
	edit3.SetLeft(220)
	edit3.SetTop(250)
	edit3.SetWidth(200)
	edit3.SetTextHint("password")

	pc1 := vcl.NewRadioButton(f)
	pc1.SetParent(f)
	pc1.SetChecked(true)
	pc1.SetCaption("PAC")
	pc1.SetLeft(150)
	pc1.SetTop(300)

	pc2 := vcl.NewRadioButton(f)
	pc2.SetParent(f)
	pc2.SetCaption("全局")
	pc2.SetLeft(250)
	pc2.SetTop(300)

	btn := vcl.NewButton(f)
	btn.SetParent(f)
	btn.SetBounds(230, 350, 90, 30)
	btn.SetCaption("action")
	btn.SetOnClick(func(sender vcl.IObject) {
		log.Println(pc1.Checked())
		log.Println(pc2.Checked())
		log.Println(edit.Text())
	})

	l2 := vcl.NewLabel(f)
	l2.SetCaption("v0.1")
	l2.SetAlign(types.AlBottom)
	l2.SetColor(colors.ClBlack)
	l2.SetParent(f)
}
