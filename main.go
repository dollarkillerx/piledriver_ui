package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/dollarkillerx/processes"
	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
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
		killPileDriver()
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
		killPileDriver()
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

	port := vcl.NewLabel(f)
	port.SetCaption("本地Socks5监听端口: ")
	port.SetBounds(60, 400, 32, 32)
	port.SetParent(f)

	edit := vcl.NewEdit(f)
	edit.SetParent(f)
	edit.SetLeft(220)
	edit.SetTop(150)
	edit.SetWidth(200)
	edit.SetTextHint("address")

	edit2 := vcl.NewEdit(f)
	edit2.SetParent(f)
	edit2.SetLeft(220)
	edit2.SetTop(200)
	edit2.SetWidth(200)
	edit2.SetTextHint("user id")

	edit3 := vcl.NewEdit(f)
	edit3.SetParent(f)
	edit3.SetLeft(220)
	edit3.SetTop(250)
	edit3.SetWidth(200)
	edit3.SetTextHint("password")

	edit4 := vcl.NewEdit(f)
	edit4.SetParent(f)
	edit4.SetLeft(220)
	edit4.SetTop(400)
	edit4.SetWidth(200)
	edit4.SetTextHint("local socks5 addr")
	edit4.SetText("0.0.0.0:8081")

	// init config
	file, err := ioutil.ReadFile("piledriver.conf")
	if err == nil {
		cfg := config{}
		if err := json.Unmarshal(file, &cfg); err == nil {
			edit.SetText(cfg.Address)
			edit2.SetText(cfg.UserID)
			edit3.SetText(cfg.Password)
			edit4.SetText(cfg.Socks5Addr)
		}
	}
	// config end

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
		if edit.Text() == "" || edit2.Text() == "" || edit3.Text() == "" || edit4.Text() == "" {
			vcl.ShowMessage("What fuck?")
			return
		}
		// 杀死旧进程
		killPileDriver()
		// run new process
		newPileDriver(edit.Text(), edit4.Text(), edit2.Text(), edit3.Text(), pc1.Checked())

		// write config
		cfg := config{Address: edit.Text(), UserID: edit2.Text(), Password: edit3.Text(), Socks5Addr: edit4.Text()}
		marshal, err := json.Marshal(cfg)
		if err == nil {
			ioutil.WriteFile("piledriver.conf", marshal, 00666)
		}

		if checkRun() {
			vcl.ShowMessage("成功启动 如果无法链接说明配置有误 或网络被拦截")
		}else {
			vcl.ShowMessage("执行失败 请检测是否配置正确")
		}
	})

	l2 := vcl.NewLabel(f)
	l2.SetCaption("v0.1")
	l2.SetAlign(types.AlBottom)
	l2.SetParent(f)
}

func checkRun() bool {
	switch runtime.GOOS {
	case "windows":
		getPid, err := processes.GetPid("piledriver_core.exe")
		if err != nil {
			return false
		}
		if getPid != "" {
			return true
		}
	case "linux":
		getPid, err := processes.GetPid("piledriver_core_linux")
		if err != nil {
			return false
		}
		if getPid != "" {
			return true
		}
	case "darwin":
		getPid, err := processes.GetPid("piledriver_core_darwin")
		if err != nil {
			return false
		}
		if getPid != "" {
			return true
		}
	}
	return false
}

func killPileDriver() error {
	var pid string

	switch runtime.GOOS {
	case "windows":
		getPid, err := processes.GetPid("piledriver_core.exe")
		if err != nil {
			return err
		}
		pid = getPid
	case "linux":
		getPid, err := processes.GetPid("piledriver_core_linux")
		if err != nil {
			return err
		}
		pid = getPid
	case "darwin":
		getPid, err := processes.GetPid("piledriver_core_darwin")
		if err != nil {
			return err
		}
		pid = getPid
	}

	return processes.KillByPid(pid)
}

func newPileDriver(addr, socks5, user, password string, pac bool) {
	switch runtime.GOOS {
	case "windows":
		if pac {
			processes.Command("./core/piledriver_core.exe", addr, socks5, user, password, "pac")
		} else {
			processes.Command("./core/piledriver_core.exe", addr, socks5, user, password)
		}
	case "linux":
		if pac {
			processes.Command("./core/piledriver_core_linux", addr, socks5, user, password, "pac")
		} else {
			processes.Command("./core/piledriver_core_linux", addr, socks5, user, password)
		}
	case "darwin":
		if pac {
			processes.Command("./core/piledriver_core_darwin", addr, socks5, user, password, "pac")
		} else {
			processes.Command("./core/piledriver_core_darwin", addr, socks5, user, password)
		}
	}
}
