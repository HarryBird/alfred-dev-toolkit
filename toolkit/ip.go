package toolkit

import (
	"log"

	"net"
	"strconv"

	alfred "github.com/HarryBird/alfred-toolkit-go"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
)

const myIpApi string = "http://ip.taobao.com/service/getIpInfo2.php"
const ipApi string = "http://ip-api.com/json/"

type myIpResp struct {
	Code int      `json:"code"`
	Data myIpData `json:"data"`
}

type myIpData struct {
	IP string `json:"ip"`
}

type ipApiResp struct {
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	ISP           string  `json:"isp"`
	ORG           string  `json:"org"`
	AS            string  `json:"as"`
	ASName        string  `json:"asname"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
}

func IPAction(ctx *cli.Context, al *alfred.Alfred) {
	args := []string(ctx.Args())
	log.Println("Args:", args)

	ip := ""

	if len(args) > 0 {
		ip = args[0]
	} else {
		ip = getMyIP()
	}

	log.Println("IP:", ip)

	if !isIPV4(ip) {
		log.Println("Invalid IP Address:" + ip)
		al.ResultAppend(alfred.NewErrorTitleItem("Invalid IP: "+ip, ""))
		al.Output()
		return
	}

	var resp ipApiResp
	url := ipApi + ip + "?fields=8114175&lang=zh-CN"

	response, _, errs := gorequest.New().Get(url).EndStruct(&resp)
	log.Println(url, response, resp, errs)

	if len(errs) > 0 {
		log.Println("IP API Fail:", url, errs)
		al.ResultAppend(alfred.NewErrorTitleItem("IP API Fail: "+ip, errs[0].Error()))
		al.Output()
		return
	}

	if response.StatusCode != 200 {
		log.Println("IP API Fail:", url, response.Status)
		al.ResultAppend(alfred.NewErrorTitleItem("IP API Fail: "+ip, response.Status))
		al.Output()
		return
	}

	if resp.Status != "success" {
		log.Println("IP API Fail:", url, response.Status)
		al.ResultAppend(alfred.NewErrorTitleItem("IP API Fail: "+ip, resp.Message))
		al.Output()
		return
	}

	al.ResultAppend(alfred.NewItem(
		"Query: "+ip,
		"",
		ip,
		ip,
		"",
		"default",
		true,
		alfred.NewDefaultIcon(),
	))

	locateSum := resp.Continent + ", " + resp.Country + ", " + resp.RegionName + ", " + resp.City + ", " + resp.District
	locateDetail := resp.Continent + "/" + resp.ContinentCode + ", " + resp.Country + "/" + resp.CountryCode + ", " + resp.RegionName + "/" + resp.Region + ", " + resp.City + ", " + resp.District

	al.ResultAppend(alfred.NewItem(
		"Location: "+locateSum,
		locateDetail,
		locateDetail,
		locateDetail,
		"",
		"default",
		true,
		alfred.NewDefaultIcon(),
	))

	geoSum := strconv.FormatFloat(resp.Lat, 'f', -1, 64) + ", " + strconv.FormatFloat(resp.Lon, 'f', -1, 64)
	geoDetail := geoSum + " " + resp.Timezone
	al.ResultAppend(alfred.NewItem(
		"GEO: "+geoSum,
		geoDetail,
		geoSum,
		geoDetail,
		"",
		"default",
		true,
		alfred.NewDefaultIcon(),
	))

	ispDetail := resp.ISP + ", " + resp.ORG
	al.ResultAppend(alfred.NewItem(
		"ISP: "+resp.ISP,
		resp.ORG,
		resp.ISP,
		ispDetail,
		"",
		"default",
		true,
		alfred.NewDefaultIcon(),
	))

	asnDetail := resp.AS + ", " + resp.ASName
	al.ResultAppend(alfred.NewItem(
		"ASN: "+resp.ASName,
		resp.AS,
		resp.AS,
		asnDetail,
		"",
		"default",
		true,
		alfred.NewDefaultIcon(),
	))

	al.ResultAppend(alfred.NewItem(
		"Mobile: "+strconv.FormatBool(resp.Mobile),
		"",
		"",
		"",
		"",
		"default",
		true,
		alfred.NewDefaultIcon(),
	))

	al.ResultAppend(alfred.NewItem(
		"Proxy: "+strconv.FormatBool(resp.Proxy),
		"",
		"",
		"",
		"",
		"default",
		true,
		alfred.NewDefaultIcon(),
	))

	al.Output()
}

func getMyIP() string {
	var resp myIpResp
	response, _, errs := gorequest.New().Post(myIpApi).Send("ip=myip").EndStruct(&resp)
	log.Println(myIpApi, response, resp, errs)

	if resp.Code == 0 {
		return resp.Data.IP
	}

	return ""
}

func isIPV4(ip string) bool {
	ipv4 := net.ParseIP(ip)

	if ipv4 == nil {
		return false
	}

	if ipv4.To4() == nil {
		return false
	}

	return true
}
