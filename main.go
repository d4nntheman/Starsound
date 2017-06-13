package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"github.com/rakyll/portmidi"
	"strconv"
	"time"
	"math/rand"
)
type Star struct {
	name string
	mag  float64
	next *Star
}
//this code is trash so ignore it//
func main() {

	//Setup midi
	portmidi.Initialize()
	fmt.Println(portmidi.CountDevices()) // returns the number of MIDI devices
	//fmt.Println(portmidi.DefaultInputDeviceID()) // returns the ID of the system default input
	//fmt.Println(portmidi.DefaultOutputDeviceID()) // returns the ID of the system default output

	//out, err := portmidi.NewOutputStream(5,1024, 0)
	out, err := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID()+2, 2048, 0)
	out2, err2 := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID()+3, 2048, 0)
	//out3, err2 := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID()+4, 2048, 0)

	if err != nil {
		log.Fatal(err)
	}
	if err2 != nil {
		log.Fatal(err)
	}
	//out.WriteShort(0xc0,0x49,0)

	//Create list head/top
	bassStars := new(Star)
	trebleStars := new(Star)

	bs := bassStars
	ts := trebleStars

	csvfile, err := os.Open("hygdata_v3.csv")
	if (err != nil) {
		log.Fatal("error opening file\n")
	}
	r := csv.NewReader(csvfile)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading all lines: %v", err)
	}

	//parse star data into LL's
	j := 0
	for i, line := range lines {
		if i == 0 {
			for p := 0; p < len(line) - 1; p++ {
				fmt.Print(p, line[p], " ")
			}
			fmt.Println()
			continue
		}
		if i == 1 {
			//fmt.Println(j, "rv:", line[11], "bf:", line[5], "prop:", line[6], "Mag:", line[13], line[33])
			mag, _ := strconv.ParseFloat(line[13], 32)//add err check
			bassStars.mag = mag
			trebleStars.mag = mag
			// i can has ?:
			if (line[6] != "") {
				bassStars.name = line[6]
				trebleStars.name = line[6]
			} else {
				bassStars.name = line[5]
				trebleStars.name = line[5]
			}
			j++

		}
		if ((line[6] != "" || line[5] != "" )) {
			//fmt.Println(j, "rv:", line[11], "bf:", line[5], "prop:", line[6], "Mag:", line[13], line[33])
			mag, _ := strconv.ParseFloat(line[13], 32)//add err check
			if (mag < 4.5 && mag > 2.5) {
				fmt.Println("bs")
				bs.next = new(Star)
				bs = bs.next
				bs.mag = mag
				if (line[6] != "") {
					bs.name = line[6]
				} else {
					bs.name = line[5]
				}
			} else if (mag >= 4.5) {
				fmt.Println("ts")
				ts.next = new(Star)
				ts = ts.next
				ts.mag = mag
				if (line[6] != "") {
					ts.name = line[6]
				} else {
					ts.name = line[5]
				}
			}
			j++
		}
	}

	//Start playing

	j = 0
	ts = trebleStars
	bs = bassStars
	ch := make(chan int64,100)
	go writesound(out2,ch)

	rand.Seed(int64(time.Now().Second()))
	rand.Seed(int64(rand.Int()))

	ran := rand.Int() % 1000
	for i:=1; i < ran; i++ {
		if(ts.next == nil){
			ts = trebleStars
		} else {
			ts = ts.next
		}
		if(bs.next == nil){
			bs = bassStars
		} else {
			bs = bs.next
		}

	}
	tonemag := 55
	//instrument settings
	out.WriteShort(0xC3,0x2C,0)
	out.WriteShort(0xC1,0x04,0)
	out.WriteShort(0xC2,0x04,0)
	out2.WriteShort(0xC1,0x04,0)

	//volume settings
	out.WriteShort(0xB3,0x07,0x30)
	for ; j < 100; {
		fmt.Print(j," ", ts.name," ", int64(10 * ts.mag), ts.mag, " ")
		j++
		if((j-1) % 13== 0){
			out.WriteShort(0x83, int64(tonemag), 100)
			tonemag+=2
			tonemag = tonemag % 66
			out.WriteShort(0x93, int64(tonemag), 60)


		}
		if (j % 4 == 0) {
			ch <- int64(bs.mag * 10)
			time.Sleep(10 * time.Millisecond)
			if (bs.next == nil) {
				bs = bassStars
			} else {
				bs = bs.next
			}
		}

		if (j % 6 == 0) {
			out.WriteShort(0x91, int64(ts.mag * 10), 100)
			out.WriteShort(0x91, int64(ts.mag * 10) + 4, 100)
			out.WriteShort(0x91, int64(ts.mag * 10) + 7, 100)

			if(rand.Int() % 100 > 30) {
				oldmag := ts.mag

				if (ts.next == nil) {
					ts = trebleStars
				} else {
					ts = ts.next
				}
				time.Sleep(333* time.Millisecond)
				out.WriteShort(0x92, int64(ts.mag * 10), 110)
				time.Sleep(333* time.Millisecond)
				out.WriteShort(0x92, int64(ts.mag * 10 + 4), 110)
				time.Sleep(333* time.Millisecond)
				out.WriteShort(0x92, int64(ts.mag * 10 + 7), 110)

				out.WriteShort(0x81, int64(oldmag * 10), 50)
				out.WriteShort(0x81, int64(oldmag * 10) + 4, 50)
				out.WriteShort(0x81, int64(oldmag * 10) + 7, 50)
				out.WriteShort(0x82, int64(ts.mag * 10), 50)
				out.WriteShort(0x82, int64(ts.mag * 10+4), 50)
				out.WriteShort(0x82, int64(ts.mag * 10+7), 50)

			} else{
				if (ts.next == nil) {
					ts = trebleStars
				} else {
					ts = ts.next
				}
				time.Sleep(333* time.Millisecond)
				out.WriteShort(0x92, int64(ts.mag * 10), 110)
				time.Sleep(333* time.Millisecond)
				out.WriteShort(0x81, int64(ts.mag * 10), 100)
				out.WriteShort(0x81, int64(ts.mag * 10) + 4, 100)
				out.WriteShort(0x81, int64(ts.mag * 10) + 7, 100)
				out.WriteShort(0x82, int64(ts.mag * 10), 50)


			}
		}

			out.WriteShort(0x92, int64(ts.mag * 10), 100)
			time.Sleep(333* time.Millisecond)
			out.WriteShort(0x82, int64(ts.mag * 10), 40)

		if (ts.next == nil) {
			ts = trebleStars
		} else {
			ts = ts.next
		}
		fmt.Println()
	}
	out.WriteShort(0x83, int64(67), 100)
	out.Close()
	out2.Close()
	portmidi.Terminate()

}

func writesound(out *portmidi.Stream, c chan int64){
	for {
		select {
		case w := <-c:
			{
				fmt.Println("bass")
				out.WriteShort(0x91, w, 100)
				out.WriteShort(0x91, w + 4, 100)
				out.WriteShort(0x91, w + 7, 100)
				time.Sleep(1000* time.Millisecond)
				out.WriteShort(0x81, w, 100)
				out.WriteShort(0x81, w + 4, 100)
				out.WriteShort(0x81, w + 7, 100)
				time.Sleep(10* time.Millisecond)
			}
		// default works here if no communication is available
		default:
		// do idle work
		}
	}
}

