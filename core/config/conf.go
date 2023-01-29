package config

import (
	"encoding/xml"
)

type Conf struct {
	XMLName xml.Name `xml:"Conf"`
	Text    string   `xml:",chardata"`
	EMM     struct {
		Text             string `xml:",chardata"`
		NetworkWhiteList struct {
			Text           string `xml:",chardata"`
			ForbidIntranet string `xml:"forbid_intranet,attr"`
			ForbidInternet string `xml:"forbid_internet,attr"`
		} `xml:"NetworkWhiteList"`
		TicketEnable    string `xml:"TicketEnable"`
		TicketLoginType string `xml:"TicketLoginType"`
		TicketLoginCode string `xml:"TicketLoginCode"`
		MdmPolicyEnable string `xml:"MdmPolicyEnable"`
		RdbUpdateTime   string `xml:"RdbUpdateTime"`
	} `xml:"EMM"`
	AworkName   string `xml:"AworkName"`
	WebSecLogin struct {
		Text              string `xml:",chardata"`
		LastLoginRes      string `xml:"LastLoginRes"`
		LastLoginTime     string `xml:"LastLoginTime"`
		LastLoginType     string `xml:"LastLoginType"`
		LastOsType        string `xml:"LastOsType"`
		LastLoginIp       string `xml:"LastLoginIp"`
		LastLoginFails    string `xml:"LastLoginFails"`
		LastLoginSuccTime string `xml:"LastLoginSuccTime"`
		LastLoginSwitch   string `xml:"LastLoginSwitch"`
	} `xml:"WebSecLogin"`
	SysTray struct {
		Text            string `xml:",chardata"`
		Enable          string `xml:"enable,attr"`
		SSLTrayIconMd5  string `xml:"SSLTrayIconMd5,attr"`
		SysShortCutName string `xml:"SysShortCutName,attr"`
		SSLTrayIconPath string `xml:"SSLTrayIconPath,attr"`
	} `xml:"SysTray"`
	Webagent struct {
		Text    string `xml:",chardata"`
		Enable  string `xml:"enable,attr"`
		Address string `xml:"address,attr"`
	} `xml:"Webagent"`
	Mline struct {
		Text     string `xml:",chardata"`
		Enable   string `xml:"enable,attr"`
		Number   string `xml:"number,attr"`
		List     string `xml:"list,attr"`
		Interval string `xml:"interval,attr"`
		Timeout  string `xml:"timeout,attr"`
	} `xml:"Mline"`
	Vpnline struct {
		Text    string `xml:",chardata"`
		Address string `xml:"address,attr"`
	} `xml:"Vpnline"`
	Htp struct {
		Text   string `xml:",chardata"`
		Enable string `xml:"enable,attr"`
		Auto   string `xml:"auto,attr"`
		Param  string `xml:"param,attr"`
		Port   string `xml:"port,attr"`
		Mtu    string `xml:"mtu,attr"`
	} `xml:"Htp"`
	WebCache struct {
		Text   string `xml:",chardata"`
		Enable string `xml:"enable,attr"`
		Mode   string `xml:"mode,attr"`
		Count  string `xml:"count,attr"`
		URL    string `xml:"url,attr"`
	} `xml:"WebCache"`
	WebOpt    string `xml:"WebOpt"`
	Bandwidth struct {
		Text      string `xml:",chardata"`
		Recvlimit string `xml:"recvlimit,attr"`
		Sendlimit string `xml:"sendlimit,attr"`
	} `xml:"Bandwidth"`
	TcpApplication struct {
		Text       string `xml:",chardata"`
		UserMode   string `xml:"userMode,attr"`
		Compress   string `xml:"compress,attr"`
		Maxthread  string `xml:"maxthread,attr"`
		Maxsession string `xml:"maxsession,attr"`
	} `xml:"TcpApplication"`
	L3VPN struct {
		Text        string `xml:",chardata"`
		IptunDns    string `xml:"iptunDns,attr"`
		IptunDnsBak string `xml:"iptunDnsBak,attr"`
	} `xml:"L3VPN"`
	SCache struct {
		Text   string `xml:",chardata"`
		Enable string `xml:"enable,attr"`
		Gwid   string `xml:"gwid,attr"`
		ID     string `xml:"id,attr"`
		Cfgmd5 string `xml:"cfgmd5,attr"`
		Dllmd5 string `xml:"dllmd5,attr"`
	} `xml:"SCache"`
	DnsRuleExceptions struct {
		Text      string   `xml:",chardata"`
		Exception []string `xml:"Exception"`
	} `xml:"DnsRuleExceptions"`
	UsbKey struct {
		Text      string `xml:",chardata"`
		Version   string `xml:"version,attr"`
		Certinput string `xml:"certinput,attr"`
		Typeinfo  string `xml:"typeinfo,attr"`
		Typecount string `xml:"typecount,attr"`
	} `xml:"UsbKey"`
	CDC struct {
		Text        string `xml:",chardata"`
		Enable      string `xml:"enable,attr"`
		LogKey      string `xml:"LogKey,attr"`
		UseUsersLog string `xml:"useUsersLog,attr"`
		LogInterval string `xml:"LogInterval,attr"`
		GrpIdInt    string `xml:"GrpIdInt,attr"`
		AuthPast    string `xml:"AuthPast,attr"`
	} `xml:"CDC"`
	Autorule struct {
		Text            string `xml:",chardata"`
		Enable          string `xml:"enable,attr"`
		EnableLimit     string `xml:"enable_limit,attr"`
		GatherRuleLimit string `xml:"gather_rule_limit,attr"`
		Mode            string `xml:"mode,attr"`
		Count           string `xml:"count,attr"`
		RuleLimit       string `xml:"rule_limit,attr"`
		Domain          string `xml:"domain,attr"`
	} `xml:"Autorule"`
	Other struct {
		Text                string `xml:",chardata"`
		LoginName           string `xml:"login_name,attr"`
		SddnEnable          string `xml:"sddn_enable,attr"`
		Sslctx              string `xml:"sslctx,attr"`
		Displayhost         string `xml:"displayhost,attr"`
		DenyNormalAccess    string `xml:"denyNormalAccess,attr"`
		EnableAutoRelogin   string `xml:"enableAutoRelogin,attr"`
		Autorelogininterval string `xml:"autorelogininterval,attr"`
		Autorelogintimes    string `xml:"autorelogintimes,attr"`
		IsRelogin           string `xml:"isRelogin,attr"`
		RelogTimeLast       string `xml:"RelogTimeLast,attr"`
		RelogIPLast         string `xml:"RelogIPLast,attr"`
		AutoStartCS         string `xml:"autoStartCS,attr"`
		PwpRemindMsg        string `xml:"pwp_remind_msg,attr"`
		Svpnlanaddr         string `xml:"svpnlanaddr,attr"`
		IsPubUser           string `xml:"isPubUser,attr"`
		IsExtern            string `xml:"isExtern,attr"`
		IsHidUser           string `xml:"isHidUser,attr"`
		Enablesavepwd       string `xml:"enablesavepwd,attr"`
		EnableCsRcWindows   string `xml:"enable_cs_rc_windows,attr"`
		Enableautologin     string `xml:"enableautologin,attr"`
		ChgPwdEnable        string `xml:"chg_pwd_enable,attr"`
		ChgPhoneEnable      string `xml:"chg_phone_enable,attr"`
		ChgNoteEnable       string `xml:"chg_note_enable,attr"`
		AuthSms             string `xml:"auth_sms,attr"`
		Mobilephone         string `xml:"mobilephone,attr"`
		UserNote            string `xml:"user_note,attr"`
		PswMinlen           string `xml:"psw_minlen,attr"`
		PptpGrpolicy        string `xml:"pptp_grpolicy,attr"`
		PptpDetaddr         string `xml:"pptp_detaddr,attr"`
		Vpntype             string `xml:"vpntype,attr"`
		Deviceversion       string `xml:"deviceversion,attr"`
		UserAuthPast        string `xml:"UserAuthPast,attr"`
		UserPwd             string `xml:"UserPwd,attr"`
		AccessibleAddr      string `xml:"AccessibleAddr,attr"`
	} `xml:"Other"`
	RemoteApp struct {
		Text              string `xml:",chardata"`
		AccountPolicy     string `xml:"account_policy,attr"`
		SessionKeeptime   string `xml:"session_keeptime,attr"`
		MapDisk           string `xml:"MapDisk,attr"`
		MapClipboard      string `xml:"MapClipboard,attr"`
		MapPrinter        string `xml:"MapPrinter,attr"`
		VirtualPrinter    string `xml:"VirtualPrinter,attr"`
		RappResuse        string `xml:"rapp_resuse,attr"`
		VirtualPrintMode  string `xml:"VirtualPrintMode,attr"`
		UseRdp            string `xml:"UseRdp,attr"`
		PrintPaper        string `xml:"PrintPaper"`
		PrivateFolderName string `xml:"PrivateFolderName"`
		SRAPOption        struct {
			Text           string `xml:",chardata"`
			LossCompressor struct {
				Text    string `xml:",chardata"`
				Type    string `xml:"type,attr"`
				Ratio   string `xml:"ratio,attr"`
				Quality string `xml:"quality,attr"`
			} `xml:"LossCompressor"`
			GlyphCompress struct {
				Text        string `xml:",chardata"`
				Option      string `xml:"option,attr"`
				JpegQuality string `xml:"jpeg_quality,attr"`
			} `xml:"GlyphCompress"`
			NoLossCompressor struct {
				Text           string `xml:",chardata"`
				BmpCompressor  string `xml:"bmp_compressor,attr"`
				CompressorType string `xml:"compressor_type,attr"`
			} `xml:"NoLossCompressor"`
			CacheHash struct {
				Text   string `xml:",chardata"`
				OpType string `xml:"op_type,attr"`
			} `xml:"CacheHash"`
			StreamMerge struct {
				Text      string `xml:",chardata"`
				Type      string `xml:"type,attr"`
				Threshold string `xml:"threshold,attr"`
				Uptime    string `xml:"uptime,attr"`
			} `xml:"StreamMerge"`
		} `xml:"SRAPOption"`
	} `xml:"RemoteApp"`
	SSLCipherSuite struct {
		Text  string `xml:",chardata"`
		EC    string `xml:"EC"`
		TCP   string `xml:"TCP"`
		L3VPN string `xml:"L3VPN"`
	} `xml:"SSLCipherSuite"`
	SSLEigenvalue struct {
		Text  string `xml:",chardata"`
		TCP   string `xml:"TCP"`
		L3VPN string `xml:"L3VPN"`
	} `xml:"SSLEigenvalue"`
	Logo struct {
		Text     string `xml:",chardata"`
		Custom   string `xml:"custom,attr"`
		LogoMd5  string `xml:"LogoMd5,attr"`
		LogoPath string `xml:"LogoPath,attr"`
	} `xml:"Logo"`
	WebHttpEnable struct {
		Text     string `xml:",chardata"`
		HttpPort string `xml:"httpPort,attr"`
		Enable   string `xml:"enable,attr"`
	} `xml:"WebHttpEnable"`
}
