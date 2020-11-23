package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"
)

const (
	frameSize = 960 * 720 * 3
)

var img = gocv.NewMat()
var classifier = gocv.NewCascadeClassifier()
var rect = image.Rectangle{}

var checker = false

func main() {
	drone := tello.NewDriver("8890")
	//window := opencv.NewWindowDriver()
	window := gocv.NewWindow("Demo2")
	classifier := gocv.NewCascadeClassifier()
	classifier.Load("haarcascade_frontalface_default.xml")
	defer classifier.Close()
	defer drone.Halt()
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
		"-an", "-sn", "-s", "960x720", "-f", "rawvideo", "pipe:1")
	ffmpegIn, _ := ffmpeg.StdinPipe()
	ffmpegOut, _ := ffmpeg.StdoutPipe()

	work := func() {
		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}
		//count:=0
		go func() {

		}()

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")

			drone.TakeOff()
			drone.StartVideo()
			drone.SetExposure(1)
			drone.SetVideoEncoderRate(4)

			go searchFace(drone)

			if checker == true {
				go follow(drone)
			}

			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start(false)
	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}

		img, err := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
		if err != nil {
			log.Print(err)
			continue
		}
		if img.Empty() {
			continue
		}
		imageRectangles := classifier.DetectMultiScale(img)

		for _, rect := range imageRectangles {
			log.Println("found a face,", rect)
			gocv.Rectangle(&img, rect, colornames.Cadetblue, 3)
		}

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
		//if count < 1000{
		//
		//}
		//window.WaitKey(1)
		//count +=1
	}
}

func searchFace(drone *tello.Driver) { //If no face found - rotate
	for {
		time.After(2 * time.Second)
		if rect.Empty() {
			drone.Clockwise(5)
		} else {
			checker = true
		}
	}
}

func follow(drone *tello.Driver) {
	for {
		if isCentered(drone) == "TOOFARLEFT" {
			drone.Right(10)
			fmt.Printf("Adjusting Right...")

		} else if isCentered(drone) == "TOOFARRIGHT" {
			drone.Left(10)
			fmt.Printf("Adjusting Left...")

		} else if isCentered(drone) == "CENTER" {
			if rect.Dx() < 150 {
				drone.Forward(5)
				fmt.Printf("Hovering Forward...")

			} else if rect.Dx() > 250 {
				drone.Backward(5)
				fmt.Printf("Hovering Backward...")

			} else {
				print("In safe range")
			}
		}
	}
}

func isCentered(drone *tello.Driver) string {
	if 960-rect.Max.X > rect.Min.X+50 {
		return "TOOFARLEFT"

	} else if 960-rect.Max.X < rect.Min.X-50 {
		return "TOOFARRIGHT"

	} else {
		return "CENTER"

	}
}
