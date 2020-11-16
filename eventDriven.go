package main

/*
import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)


import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"os/exec"
	"time"
)

func main() {
	drone := tello.NewDriver("8890")
	work := func() {
		go startCamera(drone)
		go fly(drone)
	}
	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)
	robot.Start()
}

func startCamera(drone *tello.Driver) {
	mplayer := exec.Command("mplayer", "-fps", "25", "-")
	mplayerIn, _ := mplayer.StdinPipe()
	if err := mplayer.Start(); err != nil {
		fmt.Println(err)
		return
	}
	drone.On(tello.ConnectedEvent, func(data interface{}) {
		fmt.Println("Connected")
		drone.StartVideo()
		drone.SetVideoEncoderRate(4)
		gobot.Every(100*time.Millisecond, func() {
			drone.StartVideo()
		})
	})
	drone.On(tello.VideoFrameEvent, func(data interface{}) {
		pkt := data.([]byte)
		if _, err := mplayerIn.Write(pkt); err != nil {
			fmt.Println(err)
		}
	})
}

func fly(drone *tello.Driver) {
	drone.TakeOff()
	time.Sleep(time.Second * 10)
	drone.Land()
}
*/
