package toolkit

import (
	//"fmt"
	"log"
	"strconv"
	"time"

	alfred "github.com/HarryBird/alfred-toolkit-go"
	"github.com/pkg/errors"
	ping "github.com/sparrc/go-ping"
	"github.com/urfave/cli"
)

func PingAction(ctx *cli.Context, al *alfred.Alfred) {
	args := []string(ctx.Args())
	log.Println("Args:", args)

	if len(args) == 0 {
		al.ResultAppend(alfred.NewErrorTitleItem("Invalid Host", "")).Output()
		return
	}

	host := args[0]

	num := 1

	stats, err := dail(host, num)
	if err != nil {
		al.ResultAppend(alfred.NewErrorTitleItem("Resolve Host Fail", "Host: "+host)).Output()
	} else {
		log.Println("Ping Result:", stats)
		sent := strconv.Itoa(stats.PacketsSent)
		recv := strconv.Itoa(stats.PacketsRecv)
		loss := strconv.FormatFloat(stats.PacketLoss, 'f', -1, 64)
		ip := stats.IPAddr.String()
		title := ip
		rtt := "0 ms"

		if len(stats.Rtts) > 0 {
			rtt = stats.Rtts[0].String()
		} else {
			title = ip + "(Timeout)"
		}

		//min := stats.MinRtt.String()
		//max := stats.MaxRtt.String()
		//avg := stats.AvgRtt.String()

		al.ResultAppend(buildIPItem(title, "Host: "+host+"; Sent: "+sent+"; Recv: "+recv+"; Loss: "+loss+"%; RTT: "+rtt, ip)).Output()
	}
}

func buildIPItem(title, subTitle, arg string) alfred.Item {
	return alfred.NewItem(title, subTitle, arg, arg, "", "", true, alfred.NewIcon("", "./icons/ip/ip.png"))
}

func dail(host string, n int) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(host)

	if err != nil {
		return nil, errors.Wrap(err, sign("Resolve Host Fail"))
	}

	pinger.Count = n
	pinger.Timeout = 500 * time.Millisecond
	pinger.Run()
	return pinger.Statistics(), nil
}
