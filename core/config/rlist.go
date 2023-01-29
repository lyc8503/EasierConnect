package config

import "encoding/xml"

type Resource struct {
	XMLName xml.Name `xml:"Resource"`
	Text    string   `xml:",chardata"`
	Rcs     struct {
		Text string `xml:",chardata"`
		Rc   []struct {
			Text           string `xml:",chardata"`
			ID             string `xml:"id,attr"`
			Name           string `xml:"name,attr"`
			Type           string `xml:"type,attr"`
			Proto          string `xml:"proto,attr"`
			Svc            string `xml:"svc,attr"`
			Host           string `xml:"host,attr"`
			Port           string `xml:"port,attr"`
			EnableDisguise string `xml:"enable_disguise,attr"`
			Note           string `xml:"note,attr"`
			Attr           string `xml:"attr,attr"`
			AppPath        string `xml:"app_path,attr"`
			RcGrpID        string `xml:"rc_grp_id,attr"`
			RcLogo         string `xml:"rc_logo,attr"`
			Authorization  string `xml:"authorization,attr"`
			AuthSpID       string `xml:"auth_sp_id,attr"`
			Selectid       string `xml:"selectid,attr"`
		} `xml:"Rc"`
	} `xml:"Rcs"`
	RcGroups struct {
		Text  string `xml:",chardata"`
		Group []struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"id,attr"`
			Name        string `xml:"name,attr"`
			Type        string `xml:"type,attr"`
			Logowidth   string `xml:"logowidth,attr"`
			Logoheight  string `xml:"logoheight,attr"`
			LoadBalance string `xml:"load_balance,attr"`
			ShowNote    string `xml:"show_note,attr"`
		} `xml:"Group"`
	} `xml:"RcGroups"`
	SD struct {
		Text   string `xml:",chardata"`
		Global struct {
			Text           string `xml:",chardata"`
			Enable         string `xml:"enable,attr"`
			SDRedirectFile string `xml:"SDRedirectFile,attr"`
		} `xml:"Global"`
		Policy struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
		} `xml:"Policy"`
		DesktopFormat struct {
			Text                   string `xml:",chardata"`
			Safedesk               string `xml:"safedesk,attr"`
			Com                    string `xml:"com,attr"`
			Infrared               string `xml:"infrared,attr"`
			Bluetooth              string `xml:"bluetooth,attr"`
			Printer                string `xml:"printer,attr"`
			Changedesk             string `xml:"changedesk,attr"`
			RegisterUnMon          string `xml:"register_un_mon,attr"`
			SafedeskLocalTransport string `xml:"safedesk_local_transport,attr"`
			RappInSafedesk         string `xml:"rapp_in_safedesk,attr"`
		} `xml:"DesktopFormat"`
		Internet struct {
			Text    string `xml:",chardata"`
			Tempbuf string `xml:"tempbuf,attr"`
			History string `xml:"history,attr"`
			Tables  string `xml:"tables,attr"`
			Cookies string `xml:"cookies,attr"`
		} `xml:"Internet"`
		Iplist struct {
			Text   string `xml:",chardata"`
			Iplist string `xml:"iplist,attr"`
		} `xml:"iplist"`
		Rclist struct {
			Text   string `xml:",chardata"`
			Rclist string `xml:"rclist,attr"`
		} `xml:"rclist"`
	} `xml:"SD"`
	Dns struct {
		Text      string `xml:",chardata"`
		Dnsserver string `xml:"dnsserver,attr"`
		Data      string `xml:"data,attr"`
		Filter    string `xml:"filter,attr"`
	} `xml:"Dns"`
	FileLock struct {
		Text         string `xml:",chardata"`
		Data         string `xml:"data,attr"`
		Filecount    string `xml:"filecount,attr"`
		Maxfilecount string `xml:"maxfilecount,attr"`
	} `xml:"FileLock"`
	UB struct {
		Text       string `xml:",chardata"`
		IndexInner string `xml:"index_inner,attr"`
		Ubdllinfo  string `xml:"ubdllinfo,attr"`
	} `xml:"UB"`
	Easylink struct {
		Text   string `xml:",chardata"`
		ElnkRc struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"Id"`
			ElnkRewrite string `xml:"ElnkRewrite"`
			Mode        string `xml:"Mode"`
			MapAddr     string `xml:"MapAddr"`
		} `xml:"ElnkRc"`
	} `xml:"Easylink"`
	Other struct {
		Text        string `xml:",chardata"`
		DefaultRcId string `xml:"defaultRcId,attr"`
		AllocateVip string `xml:"allocateVip,attr"`
		Balanceinfo string `xml:"balanceinfo,attr"`
	} `xml:"Other"`
	UrlWarrentRules struct {
		Text   string `xml:",chardata"`
		Enable string `xml:"enable,attr"`
		Filter string `xml:"filter,attr"`
		Tips   string `xml:"tips,attr"`
	} `xml:"UrlWarrentRules"`
	MSGINFO     string `xml:"MSG_INFO"`
	WebSsoInfos string `xml:"WebSsoInfos"`
	VSP         struct {
		Text string `xml:",chardata"`
		Misc struct {
			Text                     string `xml:",chardata"`
			SDTitle                  string `xml:"SDTitle,attr"`
			ShowRcInSD               string `xml:"ShowRcInSD,attr"`
			ShowUserShortCutIconInSD string `xml:"ShowUserShortCutIconInSD,attr"`
		} `xml:"Misc"`
		WallPaper struct {
			Text     string `xml:",chardata"`
			Type     string `xml:"Type,attr"`
			URL      string `xml:"Url,attr"`
			Compress string `xml:"Compress,attr"`
			MD5      string `xml:"MD5,attr"`
		} `xml:"WallPaper"`
		Inject struct {
			Text string `xml:",chardata"`
			Type string `xml:"Type,attr"`
		} `xml:"Inject"`
		NavBar struct {
			Text    string `xml:",chardata"`
			IconUrl string `xml:"IconUrl,attr"`
			MD5     string `xml:"MD5,attr"`
		} `xml:"NavBar"`
		OfflineVisit struct {
			Text      string `xml:",chardata"`
			Enable    string `xml:"Enable,attr"`
			VisitTime string `xml:"VisitTime,attr"`
			IsBind    string `xml:"IsBind,attr"`
		} `xml:"OfflineVisit"`
		RedirectData struct {
			Text            string `xml:",chardata"`
			ProcessType     string `xml:"ProcessType,attr"`
			UseCustomDefine string `xml:"UseCustomDefine,attr"`
		} `xml:"RedirectData"`
		Crypto struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"Type,attr"`
			Length string `xml:"Length,attr"`
			Ctx    string `xml:"Ctx,attr"`
		} `xml:"Crypto"`
		ReDirect struct {
			Text     string `xml:",chardata"`
			NameRule string `xml:"NameRule,attr"`
			Ctx      string `xml:"Ctx,attr"`
		} `xml:"ReDirect"`
		FileExport struct {
			Text             string `xml:",chardata"`
			Enable           string `xml:"Enable,attr"`
			AuditLog         string `xml:"AuditLog,attr"`
			MaxAuditFileSize string `xml:"MaxAuditFileSize,attr"`
			Compress         string `xml:"Compress,attr"`
		} `xml:"FileExport"`
		ActiveXProxyInstall struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"Enable,attr"`
		} `xml:"ActiveXProxyInstall"`
		LocalCommunication struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"Enable,attr"`
		} `xml:"LocalCommunication"`
		ExecutableProcess struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"Enable,attr"`
		} `xml:"ExecutableProcess"`
	} `xml:"VSP"`
	StaticSd struct {
		Text        string `xml:",chardata"`
		SpecileFile string `xml:"SpecileFile"`
		SpecileKey  struct {
			Text    string   `xml:",chardata"`
			KeyList []string `xml:"KeyList"`
		} `xml:"SpecileKey"`
		SpecileProc struct {
			Text     string `xml:",chardata"`
			ProcList string `xml:"ProcList"`
		} `xml:"SpecileProc"`
		DefaultExecutableProcess struct {
			Text          string `xml:",chardata"`
			WhiteListItem []struct {
				Text       string `xml:",chardata"`
				FileName   string `xml:"FileName"`
				Value      string `xml:"Value"`
				VerifyType string `xml:"VerifyType"`
			} `xml:"WhiteListItem"`
		} `xml:"DefaultExecutableProcess"`
		NotifySize          string `xml:"NotifySize"`
		EnTopTool           string `xml:"EnTopTool"`
		ActiveXProxyProcess struct {
			Text    string `xml:",chardata"`
			Process struct {
				Text string `xml:",chardata"`
				Name string `xml:"Name,attr"`
			} `xml:"Process"`
		} `xml:"ActiveXProxyProcess"`
		RedirectObjectRule struct {
			Text  string `xml:",chardata"`
			Count string `xml:"Count,attr"`
			Rule  []struct {
				Text        string `xml:",chardata"`
				ObjectType  string `xml:"ObjectType,attr"`
				Disable     string `xml:"Disable,attr"`
				MatchRule   string `xml:"MatchRule,attr"`
				ObjectName  string `xml:"ObjectName,attr"`
				ProcessName string `xml:"ProcessName,attr"`
			} `xml:"Rule"`
		} `xml:"RedirectObjectRule"`
		DeniService struct {
			Text    string `xml:",chardata"`
			Count   string `xml:"Count,attr"`
			Service []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"Service"`
		} `xml:"DeniService"`
		InterceptSet struct {
			Text  string `xml:",chardata"`
			Count string `xml:"Count,attr"`
			Item  []struct {
				Text string `xml:",chardata"`
				Type string `xml:"Type,attr"`
				Name string `xml:"Name,attr"`
			} `xml:"Item"`
		} `xml:"InterceptSet"`
		InjectAgentWhiteList struct {
			Text     string `xml:",chardata"`
			Count    string `xml:"Count,attr"`
			ProcName []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"ProcName"`
		} `xml:"InjectAgentWhiteList"`
		VBRule struct {
			Text     string   `xml:",chardata"`
			Enable   string   `xml:"Enable,attr"`
			ProcName []string `xml:"ProcName"`
		} `xml:"VBRule"`
		NetDrvInfo struct {
			Text          string `xml:",chardata"`
			WorkMode      string `xml:"WorkMode,attr"`
			WhiteListItem []struct {
				Text       string `xml:",chardata"`
				FileName   string `xml:"FileName,attr"`
				VerifyType string `xml:"VerifyType,attr"`
				Value      string `xml:"Value,attr"`
			} `xml:"WhiteListItem"`
		} `xml:"NetDrvInfo"`
		DenyProcess struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text     string `xml:",chardata"`
				FileName string `xml:"FileName,attr"`
				Info     string `xml:"Info,attr"`
			} `xml:"Item"`
		} `xml:"DenyProcess"`
		WhitePipeOfProcess struct {
			Text           string `xml:",chardata"`
			EnablePipeRule string `xml:"EnablePipeRule,attr"`
			Item           []struct {
				Text     string `xml:",chardata"`
				FileName string `xml:"FileName,attr"`
				Info     string `xml:"Info,attr"`
				PipeName string `xml:"PipeName,attr"`
			} `xml:"Item"`
		} `xml:"WhitePipeOfProcess"`
		Drivers struct {
			Text   string `xml:",chardata"`
			Enable string `xml:"Enable,attr"`
			Driver []struct {
				Text   string `xml:",chardata"`
				Name   string `xml:"Name,attr"`
				Enable string `xml:"Enable,attr"`
			} `xml:"Driver"`
		} `xml:"Drivers"`
		UsbWhiteProcess struct {
			Text string `xml:",chardata"`
			Rule []struct {
				Text        string `xml:",chardata"`
				ProcessName string `xml:"ProcessName,attr"`
				Info        string `xml:"Info,attr"`
				Type        string `xml:"Type,attr"`
			} `xml:"Rule"`
		} `xml:"UsbWhiteProcess"`
	} `xml:"StaticSd"`
}
